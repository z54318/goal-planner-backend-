package task

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"goal-planner/internal/common/response"
)

// Handler 负责处理 task 模块的 HTTP 请求。
type Handler struct {
	repo *Repository
}

// NewHandler 创建任务模块处理器。
func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		repo: NewRepository(db),
	}
}

// RegisterProtectedRoutes 注册任务模块受保护路由。
func (h *Handler) RegisterProtectedRoutes(router *gin.RouterGroup) {
	router.GET("/tasks", h.ListTasks)
	router.GET("/tasks/:id", h.GetTask)
	router.POST("/tasks", h.CreateTask)
	router.PUT("/tasks/:id", h.UpdateTask)
	router.PATCH("/tasks/:id/status", h.UpdateTaskStatus)
	router.DELETE("/tasks/:id", h.DeleteTask)
}

// ListTasks 获取任务列表
// @Summary 获取任务列表
// @Tags tasks
// @ID tasksList
// @Produce json
// @Param status query string false "任务状态"
// @Param goal_id query int false "目标ID"
// @Param phase_id query int false "阶段ID"
// @Param page query int false "页码，从1开始"
// @Param page_size query int false "每页条数"
// @Security BearerAuth
// @Success 200 {object} TaskListResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 400 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/tasks [get]
// ListTasks 返回当前登录用户的任务列表。
func (h *Handler) ListTasks(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	var req ListTasksRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "查询参数格式不正确")
		return
	}

	req.Status = TaskStatus(strings.TrimSpace(string(req.Status)))
	if req.Status != "" && !isValidTaskStatus(req.Status) {
		response.Fail(c, http.StatusBadRequest, "任务状态不合法")
		return
	}
	if req.GoalID < 0 || req.PhaseID < 0 {
		response.Fail(c, http.StatusBadRequest, "筛选参数不合法")
		return
	}
	normalizeTaskPagination(&req)

	tasks, total, err := h.repo.ListByUserID(c.Request.Context(), userID, req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询任务列表失败")
		return
	}

	response.Success(c, TaskListData{
		List:     tasks,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

// GetTask 获取任务详情
// @Summary 获取任务详情
// @Tags tasks
// @ID taskGet
// @Produce json
// @Param id path int true "任务ID"
// @Security BearerAuth
// @Success 200 {object} TaskResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/tasks/{id} [get]
// GetTask 返回当前登录用户的一条任务详情。
func (h *Handler) GetTask(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	taskID, ok := currentTaskID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "任务ID不合法")
		return
	}

	task, err := h.repo.GetByID(c.Request.Context(), userID, taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "任务不存在")
			return
		}

		response.Fail(c, http.StatusInternalServerError, "查询任务失败")
		return
	}

	response.Success(c, task)
}

// CreateTask 新增任务
// @Summary 新增任务
// @Tags tasks
// @ID taskCreate
// @Accept json
// @Produce json
// @Param request body CreateTaskRequest true "任务信息"
// @Security BearerAuth
// @Success 200 {object} TaskResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/tasks [post]
// CreateTask 为当前登录用户的阶段新增任务。
func (h *Handler) CreateTask(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求体格式不正确")
		return
	}

	if !normalizeCreateTaskRequest(&req) {
		response.Fail(c, http.StatusBadRequest, "任务参数不合法")
		return
	}

	task, err := h.repo.Create(c.Request.Context(), userID, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "阶段不存在")
			return
		}

		response.Fail(c, http.StatusInternalServerError, "创建任务失败")
		return
	}

	response.Success(c, task)
}

// UpdateTask 编辑任务
// @Summary 编辑任务
// @Tags tasks
// @ID taskUpdate
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Param request body UpdateTaskRequest true "任务信息"
// @Security BearerAuth
// @Success 200 {object} TaskResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/tasks/{id} [put]
// UpdateTask 编辑当前登录用户的一条任务。
func (h *Handler) UpdateTask(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	taskID, ok := currentTaskID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "任务ID不合法")
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求体格式不正确")
		return
	}

	if !normalizeUpdateTaskRequest(&req) {
		response.Fail(c, http.StatusBadRequest, "任务参数不合法")
		return
	}

	task, err := h.repo.Update(c.Request.Context(), userID, taskID, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "任务或阶段不存在")
			return
		}

		response.Fail(c, http.StatusInternalServerError, "更新任务失败")
		return
	}

	response.Success(c, task)
}

// UpdateTaskStatus 更新任务状态
// @Summary 更新任务状态
// @Tags tasks
// @ID taskUpdateStatus
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Param request body UpdateTaskStatusRequest true "任务状态"
// @Security BearerAuth
// @Success 200 {object} TaskResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/tasks/{id}/status [patch]
// UpdateTaskStatus 更新当前登录用户的一条任务状态。
func (h *Handler) UpdateTaskStatus(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	taskID, ok := currentTaskID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "任务ID不合法")
		return
	}

	var req UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求体格式不正确")
		return
	}

	req.Status = TaskStatus(strings.TrimSpace(string(req.Status)))
	if !isValidTaskStatus(req.Status) {
		response.Fail(c, http.StatusBadRequest, "任务状态不合法")
		return
	}

	task, err := h.repo.UpdateStatus(c.Request.Context(), userID, taskID, req.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "任务不存在")
			return
		}

		response.Fail(c, http.StatusInternalServerError, "更新任务状态失败")
		return
	}

	response.Success(c, task)
}

// DeleteTask 删除任务
// @Summary 删除任务
// @Tags tasks
// @ID taskDelete
// @Produce json
// @Param id path int true "任务ID"
// @Security BearerAuth
// @Success 200 {object} DeleteTaskResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 400 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/tasks/{id} [delete]
// DeleteTask 删除当前登录用户的一条任务。
func (h *Handler) DeleteTask(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	taskID, ok := currentTaskID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "任务ID不合法")
		return
	}

	if err := h.repo.Delete(c.Request.Context(), userID, taskID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "任务不存在")
			return
		}

		response.Fail(c, http.StatusInternalServerError, "删除任务失败")
		return
	}

	response.Success(c, DeleteTaskData{Deleted: true})
}

func isValidTaskStatus(status TaskStatus) bool {
	switch status {
	case TaskStatusTodo, TaskStatusInProgress, TaskStatusDone:
		return true
	default:
		return false
	}
}

func isValidTaskPriority(priority TaskPriority) bool {
	switch priority {
	case TaskPriorityHigh, TaskPriorityMedium, TaskPriorityLow:
		return true
	default:
		return false
	}
}

func normalizeCreateTaskRequest(req *CreateTaskRequest) bool {
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	req.Deliverables = strings.TrimSpace(req.Deliverables)
	req.Priority = TaskPriority(strings.TrimSpace(string(req.Priority)))

	if req.PhaseID <= 0 || req.Title == "" || req.EstimatedDays < 0 || req.SortOrder < 0 {
		return false
	}
	if req.Priority == "" {
		req.Priority = TaskPriorityMedium
	}

	return isValidTaskPriority(req.Priority)
}

func normalizeUpdateTaskRequest(req *UpdateTaskRequest) bool {
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	req.Deliverables = strings.TrimSpace(req.Deliverables)
	req.Priority = TaskPriority(strings.TrimSpace(string(req.Priority)))

	if req.PhaseID <= 0 || req.Title == "" || req.EstimatedDays < 0 || req.SortOrder < 0 {
		return false
	}
	if req.Priority == "" {
		req.Priority = TaskPriorityMedium
	}

	return isValidTaskPriority(req.Priority)
}

func normalizeTaskPagination(req *ListTasksRequest) {
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

func currentTaskID(c *gin.Context) (int64, bool) {
	taskIDText := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDText, 10, 64)
	if err != nil || taskID <= 0 {
		return 0, false
	}

	return taskID, true
}
