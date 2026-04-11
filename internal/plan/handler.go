package plan

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	mySQLDriver "github.com/go-sql-driver/mysql"

	"goal-planner/internal/common/response"
	appai "goal-planner/internal/infra/ai"
)

// Handler 负责处理 plan 模块的 HTTP 请求。
type Handler struct {
	repo      *Repository
	generator *appai.Client
}

// NewHandler 创建计划模块处理器。
func NewHandler(db *sql.DB, generator *appai.Client) *Handler {
	return &Handler{
		repo:      NewRepository(db),
		generator: generator,
	}
}

// RegisterProtectedRoutes 注册计划模块受保护路由。
func (h *Handler) RegisterProtectedRoutes(router *gin.RouterGroup) {
	router.GET("/goals/:id/plan", h.GetPlan)
	router.POST("/goals/:id/generate-plan", h.GeneratePlan)
	router.POST("/goals/:id/regenerate-plan", h.RegeneratePlan)
	router.PUT("/goals/:id/plan", h.UpdatePlan)
	router.DELETE("/goals/:id/plan", h.DeletePlan)
}

// GetPlan 获取目标计划
// @Summary 获取目标计划
// @Tags plans
// @ID goalPlanGet
// @Produce json
// @Param id path int true "目标ID"
// @Security BearerAuth
// @Success 200 {object} PlanResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/goals/{id}/plan [get]
// GetPlan 返回当前登录用户某个目标的计划详情。
func (h *Handler) GetPlan(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	goalID, ok := currentGoalID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "目标ID不合法")
		return
	}

	plan, err := h.repo.GetByGoalID(c.Request.Context(), userID, goalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "计划不存在")
			return
		}

		response.Fail(c, http.StatusInternalServerError, "查询计划失败")
		return
	}

	response.Success(c, plan)
}

// GeneratePlan 生成目标计划
// @Summary 生成目标计划
// @Tags plans
// @ID goalPlanGenerate
// @Produce json
// @Param id path int true "目标ID"
// @Security BearerAuth
// @Success 200 {object} PlanResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 409 {object} response.ErrorBody
// @Failure 503 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/goals/{id}/generate-plan [post]
// GeneratePlan 为当前登录用户的目标调用 AI 生成计划。
func (h *Handler) GeneratePlan(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	goalID, ok := currentGoalID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "目标ID不合法")
		return
	}

	goal, err := h.repo.GetGoalForGeneration(c.Request.Context(), userID, goalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "目标不存在")
			return
		}

		response.Fail(c, http.StatusInternalServerError, "查询目标失败")
		return
	}

	output, err := h.generator.GeneratePlan(c.Request.Context(), goal)
	if err != nil {
		if errors.Is(err, appai.ErrNotConfigured) {
			response.Fail(c, http.StatusServiceUnavailable, "AI服务未配置")
			return
		}
		if errors.Is(err, appai.ErrInvalidResponse) {
			slog.Error("generate plan failed: invalid ai response", "goal_id", goalID, "user_id", userID, "error", err)
			response.Fail(c, http.StatusBadGateway, "AI返回结果不可解析")
			return
		}
		var requestErr *appai.RequestError
		if errors.As(err, &requestErr) {
			slog.Error("generate plan failed: ai request error", "goal_id", goalID, "user_id", userID, "status_code", requestErr.StatusCode, "body", requestErr.Body)
			switch requestErr.StatusCode {
			case http.StatusUnauthorized, http.StatusForbidden:
				response.Fail(c, http.StatusBadGateway, "AI鉴权失败，请检查密钥")
				return
			case http.StatusNotFound:
				response.Fail(c, http.StatusBadGateway, "AI模型或接口地址不正确")
				return
			case http.StatusTooManyRequests:
				response.Fail(c, http.StatusBadGateway, "AI请求过于频繁或额度不足")
				return
			default:
				if requestErr.StatusCode >= http.StatusInternalServerError {
					response.Fail(c, http.StatusBadGateway, "AI服务暂时不可用")
					return
				}
			}
		}

		slog.Error("generate plan failed", "goal_id", goalID, "user_id", userID, "error", err)
		response.Fail(c, http.StatusBadGateway, "生成计划失败")
		return
	}

	plan, err := h.repo.CreateGenerated(c.Request.Context(), userID, goalID, output)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "目标不存在")
			return
		}

		var mysqlErr *mySQLDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			response.Fail(c, http.StatusConflict, "该目标已存在计划")
			return
		}

		slog.Error("save generated plan failed", "goal_id", goalID, "user_id", userID, "error", err)
		response.Fail(c, http.StatusInternalServerError, "保存计划失败")
		return
	}

	response.Success(c, plan)
}

// RegeneratePlan 重新生成目标计划
// @Summary 重新生成目标计划
// @Tags plans
// @ID goalPlanRegenerate
// @Produce json
// @Param id path int true "目标ID"
// @Security BearerAuth
// @Success 200 {object} PlanResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 503 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/goals/{id}/regenerate-plan [post]
// RegeneratePlan 删除旧计划后重新生成。
func (h *Handler) RegeneratePlan(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	goalID, ok := currentGoalID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "目标ID不合法")
		return
	}

	goal, err := h.repo.GetGoalForGeneration(c.Request.Context(), userID, goalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "目标不存在")
			return
		}

		response.Fail(c, http.StatusInternalServerError, "查询目标失败")
		return
	}

	output, err := h.generator.GeneratePlan(c.Request.Context(), goal)
	if err != nil {
		if errors.Is(err, appai.ErrNotConfigured) {
			response.Fail(c, http.StatusServiceUnavailable, "AI服务未配置")
			return
		}
		if errors.Is(err, appai.ErrInvalidResponse) {
			slog.Error("regenerate plan failed: invalid ai response", "goal_id", goalID, "user_id", userID, "error", err)
			response.Fail(c, http.StatusBadGateway, "AI返回结果不可解析")
			return
		}
		var requestErr *appai.RequestError
		if errors.As(err, &requestErr) {
			slog.Error("regenerate plan failed: ai request error", "goal_id", goalID, "user_id", userID, "status_code", requestErr.StatusCode, "body", requestErr.Body)
			switch requestErr.StatusCode {
			case http.StatusUnauthorized, http.StatusForbidden:
				response.Fail(c, http.StatusBadGateway, "AI鉴权失败，请检查密钥")
				return
			case http.StatusNotFound:
				response.Fail(c, http.StatusBadGateway, "AI模型或接口地址不正确")
				return
			case http.StatusTooManyRequests:
				response.Fail(c, http.StatusBadGateway, "AI请求过于频繁或额度不足")
				return
			default:
				if requestErr.StatusCode >= http.StatusInternalServerError {
					response.Fail(c, http.StatusBadGateway, "AI服务暂时不可用")
					return
				}
			}
		}

		slog.Error("regenerate plan failed", "goal_id", goalID, "user_id", userID, "error", err)
		response.Fail(c, http.StatusBadGateway, "重新生成计划失败")
		return
	}

	plan, err := h.repo.RegenerateGenerated(c.Request.Context(), userID, goalID, output)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "目标不存在")
			return
		}

		slog.Error("save regenerated plan failed", "goal_id", goalID, "user_id", userID, "error", err)
		response.Fail(c, http.StatusInternalServerError, "保存计划失败")
		return
	}

	response.Success(c, plan)
}

// UpdatePlan 编辑目标计划
// @Summary 编辑目标计划
// @Tags plans
// @ID goalPlanUpdate
// @Accept json
// @Produce json
// @Param id path int true "目标ID"
// @Param request body UpdatePlanRequest true "计划参数"
// @Security BearerAuth
// @Success 200 {object} PlanResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/goals/{id}/plan [put]
// UpdatePlan 更新当前登录用户某个目标下计划的标题和概述。
func (h *Handler) UpdatePlan(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	goalID, ok := currentGoalID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "目标ID不合法")
		return
	}

	var req UpdatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求体格式不正确")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Overview = strings.TrimSpace(req.Overview)
	if req.Title == "" {
		response.Fail(c, http.StatusBadRequest, "计划标题不能为空")
		return
	}
	if req.Overview == "" {
		response.Fail(c, http.StatusBadRequest, "计划概述不能为空")
		return
	}

	plan, err := h.repo.UpdateByGoalID(c.Request.Context(), userID, goalID, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "计划不存在")
			return
		}

		response.Fail(c, http.StatusInternalServerError, "更新计划失败")
		return
	}

	response.Success(c, plan)
}

// DeletePlan 删除目标计划
// @Summary 删除目标计划
// @Tags plans
// @ID goalPlanDelete
// @Produce json
// @Param id path int true "目标ID"
// @Security BearerAuth
// @Success 200 {object} response.Body
// @Failure 401 {object} response.ErrorBody
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/goals/{id}/plan [delete]
// DeletePlan 删除当前登录用户某个目标下的计划及其关联阶段和任务
func (h *Handler) DeletePlan(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		// 401
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	goalID, ok := currentGoalID(c)
	if !ok {
		// 400
		response.Fail(c, http.StatusBadRequest, "目标ID不合法")
		return
	}

	err := h.repo.DeleteByGoalID(c.Request.Context(), userID, goalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 404
			response.Fail(c, http.StatusNotFound, "计划不存在")
			return
		}
		// 500
		response.Fail(c, http.StatusInternalServerError, "删除计划失败")
		return
	}

	response.Success(c, gin.H{
		"deleted": true,
	})
}

// currentUserID 从鉴权上下文中取出当前登录用户 ID。
func currentUserID(c *gin.Context) (int64, bool) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		return 0, false
	}

	return userID, true
}

// currentGoalID 从路由参数中取出目标 ID。
func currentGoalID(c *gin.Context) (int64, bool) {
	goalIDText := c.Param("id")
	goalID, err := strconv.ParseInt(goalIDText, 10, 64)
	if err != nil || goalID <= 0 {
		return 0, false
	}

	return goalID, true
}
