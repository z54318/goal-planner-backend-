package phase

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"goal-planner/internal/common/response"
	appai "goal-planner/internal/infra/ai"
)

// Handler 负责处理 phase 模块的 HTTP 请求。
type Handler struct {
	repo      *Repository
	generator *appai.Client
}

// NewHandler 创建阶段模块处理器。
func NewHandler(db *sql.DB, generator *appai.Client) *Handler {
	return &Handler{
		repo:      NewRepository(db),
		generator: generator,
	}
}

// RegisterProtectedRoutes 注册阶段模块受保护路由。
func (h *Handler) RegisterProtectedRoutes(router *gin.RouterGroup) {
	router.GET("/phases/:id", h.GetPhase)
	router.GET("/phases/:id/next-step", h.GetSuggestion)
	router.POST("/phases/:id/next-step", h.SuggestNextStep)
	router.PUT("/phases/:id", h.UpdatePhase)
	router.DELETE("/phases/:id", h.DeletePhase)
}

// GetPhase 获取阶段详情
// @Summary 获取阶段详情
// @Tags phases
// @ID phaseGet
// @Produce json
// @Param id path int true "阶段ID"
// @Security BearerAuth
// @Success 200 {object} PhaseResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/phases/{id} [get]
// GetPhase 返回当前登录用户的一条阶段详情。
func (h *Handler) GetPhase(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	phaseID, ok := currentPhaseID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "阶段ID不合法")
		return
	}

	phase, err := h.repo.GetByID(c.Request.Context(), userID, phaseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "阶段不存在")
			return
		}

		response.Fail(c, http.StatusInternalServerError, "查询阶段失败")
		return
	}

	response.Success(c, phase)
}

// GetSuggestion 查询阶段执行建议
// @Summary 查询阶段执行建议
// @Tags phases
// @ID phaseNextStepGet
// @Produce json
// @Param id path int true "阶段ID"
// @Security BearerAuth
// @Success 200 {object} NextStepSuggestionResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/phases/{id}/next-step [get]
// GetSuggestion 返回当前登录用户的一条已保存阶段执行建议。
func (h *Handler) GetSuggestion(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	phaseID, ok := currentPhaseID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "阶段ID不合法")
		return
	}

	suggestion, err := h.repo.GetSavedSuggestionByID(c.Request.Context(), userID, phaseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "执行建议不存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询执行建议失败")
		return
	}

	response.Success(c, suggestion)
}

// SuggestNextStep 获取阶段执行建议
// @Summary 生成阶段执行建议
// @Tags phases
// @ID phaseNextStepSuggest
// @Produce json
// @Param id path int true "阶段ID"
// @Security BearerAuth
// @Success 200 {object} NextStepSuggestionResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 502 {object} response.ErrorBody
// @Failure 503 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/phases/{id}/next-step [post]
// SuggestNextStep 为当前登录用户的一条阶段生成并保存执行建议。
func (h *Handler) SuggestNextStep(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	phaseID, ok := currentPhaseID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "阶段ID不合法")
		return
	}

	input, err := h.repo.GetSuggestionContextByID(c.Request.Context(), userID, phaseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "阶段不存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询阶段上下文失败")
		return
	}

	suggestion, err := h.generator.SuggestNextStepForPhase(c.Request.Context(), input)
	if err != nil {
		handleSuggestionError(c, err, "生成阶段执行建议失败")
		return
	}

	if err := h.repo.SaveSuggestionByID(c.Request.Context(), userID, phaseID, suggestion); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "阶段不存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "保存执行建议失败")
		return
	}

	response.Success(c, suggestion)
}

// UpdatePhase 编辑阶段
// @Summary 编辑阶段
// @Tags phases
// @ID phaseUpdate
// @Accept json
// @Produce json
// @Param id path int true "阶段ID"
// @Param request body UpdatePhaseRequest true "阶段信息"
// @Security BearerAuth
// @Success 200 {object} PhaseResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/phases/{id} [put]
// UpdatePhase 编辑当前登录用户的一条阶段。
func (h *Handler) UpdatePhase(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	phaseID, ok := currentPhaseID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "阶段ID不合法")
		return
	}

	var req UpdatePhaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求体格式不正确")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	if req.Title == "" || req.SortOrder < 0 {
		response.Fail(c, http.StatusBadRequest, "阶段参数不合法")
		return
	}

	phase, err := h.repo.Update(c.Request.Context(), userID, phaseID, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "阶段不存在")
			return
		}

		response.Fail(c, http.StatusInternalServerError, "更新阶段失败")
		return
	}

	response.Success(c, phase)
}

// DeletePhase 删除阶段
// @Summary 删除阶段
// @Tags phases
// @ID phaseDelete
// @Produce json
// @Param id path int true "阶段ID"
// @Security BearerAuth
// @Success 200 {object} response.Body
// @Failure 401 {object} response.ErrorBody
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/phases/{id} [delete]
// DeletePhase 删除当前用户的一条阶段及其关联任务。
func (h *Handler) DeletePhase(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	phaseID, ok := currentPhaseID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "阶段ID不合法")
		return
	}

	if err := h.repo.Delete(c.Request.Context(), userID, phaseID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "阶段不存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "删除阶段失败")
		return
	}
	response.Success(c, gin.H{
		"deleted": true,
	})
}

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

func currentPhaseID(c *gin.Context) (int64, bool) {
	phaseIDText := c.Param("id")
	phaseID, err := strconv.ParseInt(phaseIDText, 10, 64)
	if err != nil || phaseID <= 0 {
		return 0, false
	}

	return phaseID, true
}

func handleSuggestionError(c *gin.Context, err error, fallbackMessage string) {
	if errors.Is(err, appai.ErrNotConfigured) {
		response.Fail(c, http.StatusServiceUnavailable, "AI服务未配置")
		return
	}
	if errors.Is(err, appai.ErrInvalidResponse) {
		response.Fail(c, http.StatusBadGateway, "AI返回结果不可解析")
		return
	}

	var requestErr *appai.RequestError
	if errors.As(err, &requestErr) {
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

	response.Fail(c, http.StatusBadGateway, fallbackMessage)
}
