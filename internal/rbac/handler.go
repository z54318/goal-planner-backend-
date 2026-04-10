package rbac

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	mySQLDriver "github.com/go-sql-driver/mysql"

	"goal-planner/internal/common/response"
)

// Handler 负责处理 RBAC 管理相关请求。
type Handler struct {
	repo *Repository
}

// NewHandler 创建 RBAC 模块处理器。
func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		repo: NewRepository(db),
	}
}

// RegisterProtectedRoutes 注册 RBAC 管理受保护路由。
func (h *Handler) RegisterProtectedRoutes(router *gin.RouterGroup) {
	router.GET("/admin/roles", h.ListRoles)
	router.GET("/admin/roles/:id/permissions", h.GetRolePermissions)
	router.PUT("/admin/roles/:id/permissions", h.UpdateRolePermissions)

	router.GET("/admin/permissions", h.ListPermissions)
	router.POST("/admin/permissions", h.CreatePermission)
	router.PUT("/admin/permissions/:id", h.UpdatePermission)
	router.DELETE("/admin/permissions/:id", h.DeletePermission)
}

// ListRoles 获取角色列表
// @Summary 获取角色列表
// @Tags rbac
// @ID adminRolesList
// @Produce json
// @Security BearerAuth
// @Success 200 {object} RoleListResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/admin/roles [get]
func (h *Handler) ListRoles(c *gin.Context) {
	roles, err := h.repo.ListRoles(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询角色失败")
		return
	}

	response.Success(c, roles)
}

// GetRolePermissions 获取角色权限绑定
// @Summary 获取角色权限绑定
// @Tags rbac
// @ID adminRolePermissionsGet
// @Produce json
// @Param id path int true "角色ID"
// @Security BearerAuth
// @Success 200 {object} RolePermissionIDsResponse
// @Failure 400 {object} response.ErrorBody
// @Failure 401 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/admin/roles/{id}/permissions [get]
func (h *Handler) GetRolePermissions(c *gin.Context) {
	roleID, ok := currentID(c, "id")
	if !ok {
		response.Fail(c, http.StatusBadRequest, "角色ID不合法")
		return
	}

	permissionIDs, err := h.repo.GetRolePermissionIDs(c.Request.Context(), roleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "角色不存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询角色权限失败")
		return
	}

	response.Success(c, RolePermissionIDsData{
		RoleID:        roleID,
		PermissionIDs: permissionIDs,
	})
}

// UpdateRolePermissions 更新角色权限绑定
// @Summary 更新角色权限绑定
// @Tags rbac
// @ID adminRolePermissionsUpdate
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param request body UpdateRolePermissionsRequest true "角色权限参数"
// @Security BearerAuth
// @Success 200 {object} RolePermissionIDsResponse
// @Failure 400 {object} response.ErrorBody
// @Failure 401 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/admin/roles/{id}/permissions [put]
func (h *Handler) UpdateRolePermissions(c *gin.Context) {
	roleID, ok := currentID(c, "id")
	if !ok {
		response.Fail(c, http.StatusBadRequest, "角色ID不合法")
		return
	}

	var req UpdateRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求体格式不正确")
		return
	}

	req.PermissionIDs = normalizePermissionIDs(req.PermissionIDs)
	if err := h.repo.ReplaceRolePermissions(c.Request.Context(), roleID, req.PermissionIDs); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "角色不存在")
			return
		}
		if errors.Is(err, ErrPermissionNotFound) {
			response.Fail(c, http.StatusBadRequest, "存在无效的权限ID")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "更新角色权限失败")
		return
	}

	response.Success(c, RolePermissionIDsData{
		RoleID:        roleID,
		PermissionIDs: req.PermissionIDs,
	})
}

// ListPermissions 获取权限列表
// @Summary 获取权限列表
// @Tags rbac
// @ID adminPermissionsList
// @Produce json
// @Security BearerAuth
// @Success 200 {object} PermissionListResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/admin/permissions [get]
func (h *Handler) ListPermissions(c *gin.Context) {
	permissions, err := h.repo.ListPermissions(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询权限失败")
		return
	}

	response.Success(c, permissions)
}

// CreatePermission 新增权限
// @Summary 新增权限
// @Tags rbac
// @ID adminPermissionCreate
// @Accept json
// @Produce json
// @Param request body CreatePermissionRequest true "权限参数"
// @Security BearerAuth
// @Success 200 {object} PermissionResponse
// @Failure 400 {object} response.ErrorBody
// @Failure 401 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/admin/permissions [post]
func (h *Handler) CreatePermission(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求体格式不正确")
		return
	}

	normalizePermissionNames(&req.Name, &req.Code)
	if req.Name == "" || req.Code == "" {
		response.Fail(c, http.StatusBadRequest, "权限名称和编码不能为空")
		return
	}

	permission, err := h.repo.CreatePermission(c.Request.Context(), req)
	if err != nil {
		var mysqlErr *mySQLDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			response.Fail(c, http.StatusBadRequest, "权限编码已存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "创建权限失败")
		return
	}

	response.Success(c, permission)
}

// UpdatePermission 更新权限
// @Summary 更新权限
// @Tags rbac
// @ID adminPermissionUpdate
// @Accept json
// @Produce json
// @Param id path int true "权限ID"
// @Param request body UpdatePermissionRequest true "权限参数"
// @Security BearerAuth
// @Success 200 {object} PermissionResponse
// @Failure 400 {object} response.ErrorBody
// @Failure 401 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/admin/permissions/{id} [put]
func (h *Handler) UpdatePermission(c *gin.Context) {
	permissionID, ok := currentID(c, "id")
	if !ok {
		response.Fail(c, http.StatusBadRequest, "权限ID不合法")
		return
	}

	var req UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求体格式不正确")
		return
	}

	normalizePermissionNames(&req.Name, &req.Code)
	if req.Name == "" || req.Code == "" {
		response.Fail(c, http.StatusBadRequest, "权限名称和编码不能为空")
		return
	}

	permission, err := h.repo.UpdatePermission(c.Request.Context(), permissionID, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "权限不存在")
			return
		}
		var mysqlErr *mySQLDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			response.Fail(c, http.StatusBadRequest, "权限编码已存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "更新权限失败")
		return
	}

	response.Success(c, permission)
}

// DeletePermission 删除权限
// @Summary 删除权限
// @Tags rbac
// @ID adminPermissionDelete
// @Produce json
// @Param id path int true "权限ID"
// @Security BearerAuth
// @Success 200 {object} response.Body
// @Failure 400 {object} response.ErrorBody
// @Failure 401 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/admin/permissions/{id} [delete]
func (h *Handler) DeletePermission(c *gin.Context) {
	permissionID, ok := currentID(c, "id")
	if !ok {
		response.Fail(c, http.StatusBadRequest, "权限ID不合法")
		return
	}

	if err := h.repo.DeletePermission(c.Request.Context(), permissionID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "权限不存在")
			return
		}
		if errors.Is(err, ErrPermissionInUse) {
			response.Fail(c, http.StatusBadRequest, "权限仍被角色或菜单使用")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "删除权限失败")
		return
	}

	response.Success(c, gin.H{
		"deleted": true,
	})
}

func currentID(c *gin.Context, name string) (int64, bool) {
	text := c.Param(name)
	id, err := strconv.ParseInt(text, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

func normalizePermissionNames(name *string, code *string) {
	*name = strings.TrimSpace(*name)
	*code = strings.TrimSpace(*code)
}

func normalizePermissionIDs(ids []int64) []int64 {
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
