package phase

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"goal-planner/internal/common/response"
)

// Handler 负责处理 phase 模块的 HTTP 请求。
type Handler struct {
	repo *Repository
}

// NewHandler 创建阶段模块处理器。
func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		repo: NewRepository(db),
	}
}

// RegisterProtectedRoutes 注册阶段模块受保护路由。
func (h *Handler) RegisterProtectedRoutes(router *gin.RouterGroup) {
	router.GET("/phases/:id", h.GetPhase)
	router.PUT("/phases/:id", h.UpdatePhase)
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
