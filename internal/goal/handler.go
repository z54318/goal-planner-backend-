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

// RegisterRoutes 注册目标模块路由。
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/api/goals", h.ListGoals)
	router.GET("/api/goals/:id", h.GetGoal)
	router.POST("/api/goals", h.CreateGoal)
}

// ListGoals 获取目标列表
// @Summary 获取目标列表
// @Tags goals
// @ID goalsList
// @Produce json
// @Success 200 {object} GoalListResponse
// @Failure 500 {object} response.ErrorBody
// @Router /api/goals [get]
// ListGoals 返回数据库中的目标列表。
func (h *Handler) ListGoals(c *gin.Context) {
	goals, err := h.repo.List(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询目标列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":  goals,
		"total": len(goals),
	})
}

// GetGoal 获取目标详情
// @Summary 获取目标详情
// @Tags goals
// @ID goalGet
// @Produce json
// @Param id path int true "目标ID"
// @Success 200 {object} GoalResponse
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/goals/{id} [get]
// GetGoal 返回单个目标详情。
func (h *Handler) GetGoal(c *gin.Context) {
	// 从路由参数中读取目标 ID。
	idText := c.Param("id")
	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, http.StatusBadRequest, "目标ID不合法")
		return
	}

	goal, err := h.repo.GetByID(c.Request.Context(), id)
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
// @Success 200 {object} GoalResponse
// @Failure 400 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/goals [post]
// CreateGoal 创建一条新的目标记录。
func (h *Handler) CreateGoal(c *gin.Context) {
	var req CreateGoalRequest

	// 解析前端传来的 JSON 请求体。
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求体格式不正确")
		return
	}

	// 目标标题不能为空。
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)

	if req.Title == "" {
		response.Fail(c, http.StatusBadRequest, "目标标题不能为空")
		return
	}

	goal, err := h.repo.Create(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建目标失败")
		return
	}

	response.Success(c, goal)
}
