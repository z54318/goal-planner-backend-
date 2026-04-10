package goal

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"goal-planner/internal/common/response"
)

// Handler 负责处理 goal 模块的 HTTP 请求。
type Handler struct {
	repo *Repository
}

// NewHandler 创建目标模块处理器。
func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		repo: NewRepository(db),
	}
}

// RegisterProtectedRoutes 注册目标模块受保护路由。
func (h *Handler) RegisterProtectedRoutes(router *gin.RouterGroup) {
	router.GET("/goals", h.ListGoals)
	router.GET("/goals/:id", h.GetGoal)
	router.POST("/goals", h.CreateGoal)
	router.PUT("/goals/:id", h.UpdateGoal)
	router.PATCH("/goals/:id/status", h.UpdateGoalStatus)
	router.DELETE("/goals/:id", h.DeleteGoal)
}

// ListGoals 获取目标列表
// @Summary 获取目标列表
// @Tags goals
// @ID goalsList
// @Produce json
// @Param page query int false "页码，从1开始"
// @Param page_size query int false "每页条数"
// @Security BearerAuth
// @Success 200 {object} GoalListResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 400 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/goals [get]
// ListGoals 返回当前登录用户的目标列表。
func (h *Handler) ListGoals(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	var req ListGoalsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "查询参数格式不正确")
		return
	}

	normalizeGoalPagination(&req)

	goals, total, err := h.repo.ListByUserID(c.Request.Context(), userID, req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询目标列表失败")
		return
	}

	response.Success(c, GoalListData{
		List:     goals,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

// GetGoal 获取目标详情
// @Summary 获取目标详情
// @Tags goals
// @ID goalGet
// @Produce json
// @Param id path int true "目标ID"
// @Security BearerAuth
// @Success 200 {object} GoalResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/goals/{id} [get]
// GetGoal 返回当前登录用户的单个目标详情。
func (h *Handler) GetGoal(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	idText := c.Param("id")
	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, http.StatusBadRequest, "目标ID不合法")
		return
	}

	goal, err := h.repo.GetByID(c.Request.Context(), userID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "目标不存在")
			return
		}

		response.Fail(c, http.StatusInternalServerError, "查询目标详情失败")
		return
	}

	response.Success(c, goal)
}

// CreateGoal 创建目标
// @Summary 创建目标
// @Tags goals
// @ID goalCreate
// @Accept json
// @Produce json
// @Param request body CreateGoalRequest true "目标参数"
// @Security BearerAuth
// @Success 200 {object} GoalResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 400 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/goals [post]
// CreateGoal 创建一条新的目标记录。
func (h *Handler) CreateGoal(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	var req CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求体格式不正确")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	req.Category = strings.TrimSpace(req.Category)

	if req.Title == "" {
		response.Fail(c, http.StatusBadRequest, "目标标题不能为空")
		return
	}

	goal, err := h.repo.Create(c.Request.Context(), userID, req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建目标失败")
		return
	}

	response.Success(c, goal)
}

// UpdateGoal 更新目标
// @Summary 更新目标
// @Tags goals
// @ID goalUpdate
// @Accept json
// @Produce json
// @Param id path int true "目标ID"
// @Param request body UpdateGoalRequest true "目标参数"
// @Security BearerAuth
// @Success 200 {object} GoalResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/goals/{id} [put]
// UpdateGoal 更新一条目标记录的主要内容。
func (h *Handler) UpdateGoal(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	idText := c.Param("id")
	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, http.StatusBadRequest, "目标ID不合法")
		return
	}

	var req UpdateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求体格式不正确")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	req.Category = strings.TrimSpace(req.Category)

	if req.Title == "" {
		response.Fail(c, http.StatusBadRequest, "目标标题不能为空")
		return
	}

	goal, err := h.repo.Update(c.Request.Context(), userID, id, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "目标不存在")
			return
		}

		response.Fail(c, http.StatusInternalServerError, "更新目标失败")
		return
	}

	response.Success(c, goal)
}

// UpdateGoalStatus 更新目标状态
// @Summary 更新目标状态
// @Tags goals
// @ID goalUpdateStatus
// @Accept json
// @Produce json
// @Param id path int true "目标ID"
// @Param request body UpdateGoalStatusRequest true "目标状态参数"
// @Security BearerAuth
// @Success 200 {object} GoalResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/goals/{id}/status [patch]
// UpdateGoalStatus 只更新目标状态。
func (h *Handler) UpdateGoalStatus(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	idText := c.Param("id")
	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, http.StatusBadRequest, "目标ID不合法")
		return
	}

	var req UpdateGoalStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求体格式不正确")
		return
	}

	req.Status = GoalStatus(strings.TrimSpace(string(req.Status)))
	if !isValidGoalStatus(req.Status) {
		response.Fail(c, http.StatusBadRequest, "目标状态不合法")
		return
	}

	goal, err := h.repo.UpdateStatus(c.Request.Context(), userID, id, req.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "目标不存在")
			return
		}

		response.Fail(c, http.StatusInternalServerError, "更新目标状态失败")
		return
	}

	response.Success(c, goal)
}

// DeleteGoal 删除目标
// @Summary 删除目标
// @Tags goals
// @ID goalDelete
// @Produce json
// @Param id path int true "目标ID"
// @Security BearerAuth
// @Success 200 {object} response.Body
// @Failure 401 {object} response.ErrorBody
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/goals/{id} [delete]
// DeleteGoal 删除一条目标记录。
func (h *Handler) DeleteGoal(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	idText := c.Param("id")
	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, http.StatusBadRequest, "目标ID不合法")
		return
	}

	if err := h.repo.Delete(c.Request.Context(), userID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "目标不存在")
			return
		}

		response.Fail(c, http.StatusInternalServerError, "删除目标失败")
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

// isValidGoalStatus 校验目标状态是否合法。
func isValidGoalStatus(status GoalStatus) bool {
	switch status {
	case GoalStatusDraft, GoalStatusActive, GoalStatusCompleted, GoalStatusArchived:
		return true
	default:
		return false
	}
}

func normalizeGoalPagination(req *ListGoalsRequest) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}
}
