package user

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"goal-planner/internal/common/response"
)

// Handler 负责处理 user 模块的 HTTP 请求。
type Handler struct {
	repo *Repository
}

// NewHandler 创建用户模块处理器。
func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		repo: NewRepository(db),
	}
}

// RegisterProtectedRoutes 注册用户管理受保护路由。
func (h *Handler) RegisterProtectedRoutes(router *gin.RouterGroup) {
	router.GET("/admin/users", h.ListUsers)
	router.PUT("/admin/users/:id/roles", h.UpdateUserRoles)
}

// ListUsers 获取用户列表
// @Summary 获取用户列表
// @Tags users
// @ID adminUsersList
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserListResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/admin/users [get]
func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.repo.List(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询用户列表失败")
		return
	}

	response.Success(c, users)
}

// UpdateUserRoles 更新用户角色
// @Summary 更新用户角色
// @Tags users
// @ID adminUserRolesUpdate
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param request body UpdateUserRolesRequest true "用户角色参数"
// @Security BearerAuth
// @Success 200 {object} UserResponse
// @Failure 400 {object} response.ErrorBody
// @Failure 401 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/admin/users/{id}/roles [put]
func (h *Handler) UpdateUserRoles(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "用户ID不合法")
		return
	}

	var req UpdateUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求体格式不正确")
		return
	}

	req.RoleIDs = normalizeRoleIDs(req.RoleIDs)

	user, err := h.repo.ReplaceRoles(c.Request.Context(), userID, req.RoleIDs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "用户不存在")
			return
		}
		if errors.Is(err, ErrRoleNotFound) {
			response.Fail(c, http.StatusBadRequest, "存在无效的角色ID")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "更新用户角色失败")
		return
	}

	response.Success(c, user)
}

func currentUserID(c *gin.Context) (int64, bool) {
	text := c.Param("id")
	id, err := strconv.ParseInt(text, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

func normalizeRoleIDs(ids []int64) []int64 {
	seen := make(map[int64]struct{}, len(ids))
	result := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}
