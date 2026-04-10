package menu

import (
	"database/sql"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"goal-planner/internal/common/response"
)

// Handler 负责处理 menu 模块的 HTTP 请求。
type Handler struct {
	repo *Repository
}

// NewHandler 创建菜单模块处理器。
func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		repo: NewRepository(db),
	}
}

// RegisterProtectedRoutes 注册菜单管理受保护路由。
func (h *Handler) RegisterProtectedRoutes(router *gin.RouterGroup) {
	router.GET("/admin/menus", h.ListMenus)
	router.GET("/admin/menus/:id", h.GetMenu)
	router.POST("/admin/menus", h.CreateMenu)
	router.PUT("/admin/menus/:id", h.UpdateMenu)
	router.DELETE("/admin/menus/:id", h.DeleteMenu)
}

// ListMenus 获取菜单列表
// @Summary 获取菜单列表
// @Tags menus
// @ID adminMenusList
// @Produce json
// @Security BearerAuth
// @Success 200 {object} MenuListResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/admin/menus [get]
// ListMenus 返回后台管理菜单树。
func (h *Handler) ListMenus(c *gin.Context) {
	menus, err := h.repo.List(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询菜单失败")
		return
	}

	response.Success(c, buildMenuTree(menus))
}

// GetMenu 获取菜单详情
// @Summary 获取菜单详情
// @Tags menus
// @ID adminMenuGet
// @Produce json
// @Param id path int true "菜单ID"
// @Security BearerAuth
// @Success 200 {object} MenuResponse
// @Failure 400 {object} response.ErrorBody
// @Failure 401 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/admin/menus/{id} [get]
// GetMenu 返回单个菜单详情。
func (h *Handler) GetMenu(c *gin.Context) {
	id, ok := currentMenuID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "菜单ID不合法")
		return
	}

	menu, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "菜单不存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询菜单失败")
		return
	}

	response.Success(c, menu)
}

// CreateMenu 新增菜单
// @Summary 新增菜单
// @Tags menus
// @ID adminMenuCreate
// @Accept json
// @Produce json
// @Param request body CreateMenuRequest true "菜单参数"
// @Security BearerAuth
// @Success 200 {object} MenuResponse
// @Failure 400 {object} response.ErrorBody
// @Failure 401 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/admin/menus [post]
// CreateMenu 创建一条菜单记录。
func (h *Handler) CreateMenu(c *gin.Context) {
	var req CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求体格式不正确")
		return
	}

	if !normalizeAndValidateCreateRequest(&req) {
		response.Fail(c, http.StatusBadRequest, "菜单名称和路由路径不能为空")
		return
	}

	parentExists, err := h.repo.ParentExists(c.Request.Context(), req.ParentID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "校验父菜单失败")
		return
	}
	if !parentExists {
		response.Fail(c, http.StatusBadRequest, "父菜单不存在")
		return
	}

	menu, err := h.repo.Create(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建菜单失败")
		return
	}

	response.Success(c, menu)
}

// UpdateMenu 更新菜单
// @Summary 更新菜单
// @Tags menus
// @ID adminMenuUpdate
// @Accept json
// @Produce json
// @Param id path int true "菜单ID"
// @Param request body UpdateMenuRequest true "菜单参数"
// @Security BearerAuth
// @Success 200 {object} MenuResponse
// @Failure 400 {object} response.ErrorBody
// @Failure 401 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/admin/menus/{id} [put]
// UpdateMenu 更新一条菜单记录。
func (h *Handler) UpdateMenu(c *gin.Context) {
	id, ok := currentMenuID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "菜单ID不合法")
		return
	}

	var req UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求体格式不正确")
		return
	}

	if !normalizeAndValidateUpdateRequest(&req) {
		response.Fail(c, http.StatusBadRequest, "菜单名称和路由路径不能为空")
		return
	}
	if req.ParentID == id {
		response.Fail(c, http.StatusBadRequest, "父菜单不能是自身")
		return
	}

	parentExists, err := h.repo.ParentExists(c.Request.Context(), req.ParentID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "校验父菜单失败")
		return
	}
	if !parentExists {
		response.Fail(c, http.StatusBadRequest, "父菜单不存在")
		return
	}

	menu, err := h.repo.Update(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "菜单不存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "更新菜单失败")
		return
	}

	response.Success(c, menu)
}

// DeleteMenu 删除菜单
// @Summary 删除菜单
// @Tags menus
// @ID adminMenuDelete
// @Produce json
// @Param id path int true "菜单ID"
// @Security BearerAuth
// @Success 200 {object} response.Body
// @Failure 400 {object} response.ErrorBody
// @Failure 401 {object} response.ErrorBody
// @Failure 404 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/admin/menus/{id} [delete]
// DeleteMenu 删除一条菜单记录。
func (h *Handler) DeleteMenu(c *gin.Context) {
	id, ok := currentMenuID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "菜单ID不合法")
		return
	}

	err := h.repo.Delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "菜单不存在")
			return
		}
		if errors.Is(err, ErrMenuHasChildren) {
			response.Fail(c, http.StatusBadRequest, "请先删除子菜单")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "删除菜单失败")
		return
	}

	response.Success(c, gin.H{
		"deleted": true,
	})
}

func normalizeAndValidateCreateRequest(req *CreateMenuRequest) bool {
	req.Name = strings.TrimSpace(req.Name)
	req.Path = strings.TrimSpace(req.Path)
	req.Component = strings.TrimSpace(req.Component)
	req.Icon = strings.TrimSpace(req.Icon)
	req.PermissionCode = strings.TrimSpace(req.PermissionCode)
	return req.Name != "" && req.Path != "" && req.ParentID >= 0
}

func normalizeAndValidateUpdateRequest(req *UpdateMenuRequest) bool {
	req.Name = strings.TrimSpace(req.Name)
	req.Path = strings.TrimSpace(req.Path)
	req.Component = strings.TrimSpace(req.Component)
	req.Icon = strings.TrimSpace(req.Icon)
	req.PermissionCode = strings.TrimSpace(req.PermissionCode)
	return req.Name != "" && req.Path != "" && req.ParentID >= 0
}

func currentMenuID(c *gin.Context) (int64, bool) {
	menuIDText := c.Param("id")
	menuID, err := strconv.ParseInt(menuIDText, 10, 64)
	if err != nil || menuID <= 0 {
		return 0, false
	}
	return menuID, true
}

func buildMenuTree(menus []Menu) []Menu {
	menuMap := make(map[int64]*Menu, len(menus))
	for i := range menus {
		menus[i].Children = nil
		menuMap[menus[i].ID] = &menus[i]
	}

	rootMenus := make([]*Menu, 0)
	for i := range menus {
		menu := &menus[i]
		if menu.ParentID == 0 {
			rootMenus = append(rootMenus, menu)
			continue
		}

		parent, ok := menuMap[menu.ParentID]
		if !ok {
			rootMenus = append(rootMenus, menu)
			continue
		}

		parent.Children = append(parent.Children, *menu)
	}

	sort.Slice(rootMenus, func(i, j int) bool {
		if rootMenus[i].SortOrder == rootMenus[j].SortOrder {
			return rootMenus[i].ID < rootMenus[j].ID
		}
		return rootMenus[i].SortOrder < rootMenus[j].SortOrder
	})

	var sortChildren func(items []Menu)
	sortChildren = func(items []Menu) {
		sort.Slice(items, func(i, j int) bool {
			if items[i].SortOrder == items[j].SortOrder {
				return items[i].ID < items[j].ID
			}
			return items[i].SortOrder < items[j].SortOrder
		})
		for i := range items {
			if len(items[i].Children) > 0 {
				sortChildren(items[i].Children)
			}
		}
	}

	tree := make([]Menu, 0, len(rootMenus))
	for _, menu := range rootMenus {
		if len(menu.Children) > 0 {
			sortChildren(menu.Children)
		}
		tree = append(tree, *menu)
	}

	return tree
}
