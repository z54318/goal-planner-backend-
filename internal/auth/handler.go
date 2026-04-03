package auth

import (
	"database/sql"
	"errors"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"goal-planner/internal/common/response"
	appjwt "goal-planner/internal/infra/jwt"
)

// Handler 负责处理认证相关请求。
type Handler struct {
	repo       *Repository
	jwtManager *appjwt.Manager
}

// NewHandler 创建认证处理器。
func NewHandler(db *sql.DB, jwtManager *appjwt.Manager) *Handler {
	return &Handler{
		repo:       NewRepository(db),
		jwtManager: jwtManager,
	}
}

// RegisterRoutes 注册认证模块路由。
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/api/auth/register", h.Register)
	router.POST("/api/auth/login", h.Login)
}

// Register 用户注册
// @Summary 用户注册
// @Tags auth
// @ID authRegister
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册参数"
// @Success 200 {object} RegisterResponse
// @Failure 400 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/auth/register [post]
// Register 处理用户注册。
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	// 解析前端传来的注册请求体。
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求体格式不正确")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Username == "" || req.Email == "" || req.Password == "" {
		response.Fail(c, http.StatusBadRequest, "用户名、邮箱和密码不能为空")
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "密码加密失败")
		return
	}

	userID, err := h.repo.CreateUser(c.Request.Context(), req.Username, req.Email, string(passwordHash))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "注册失败，用户名或邮箱可能已存在")
		return
	}

	response.Success(c, gin.H{
		"user_id":  userID,
		"username": req.Username,
		"email":    req.Email,
	})
}

// Login 用户登录
// @Summary 用户登录
// @Tags auth
// @ID authLogin
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录参数"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} response.ErrorBody
// @Failure 401 {object} response.ErrorBody
// @Failure 403 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/auth/login [post]
// Login 处理用户登录。
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	// 解析前端传来的登录请求体。
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "请求体格式不正确")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	if req.Username == "" || req.Password == "" {
		response.Fail(c, http.StatusBadRequest, "用户名或密码不能为空")
		return
	}

	user, err := h.repo.GetByUsername(c.Request.Context(), req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Fail(c, http.StatusUnauthorized, "用户名或密码错误")
			return
		}

		response.Fail(c, http.StatusInternalServerError, "查询用户失败")
		return
	}

	if user.Status != "active" {
		response.Fail(c, http.StatusForbidden, "用户已被禁用")
		return
	}

	// 校验用户输入密码和数据库中的密码哈希是否匹配。
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.Fail(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	token, err := h.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "生成登录凭证失败")
		return
	}

	// 登录成功后返回 JWT，后续前端通过 Authorization 请求头携带。
	response.Success(c, gin.H{
		"token":    token,
		"user_id":  user.ID,
		"username": user.Username,
	})
}

// Profile 当前用户信息
// @Summary 获取当前登录用户信息
// @Tags auth
// @ID authProfile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ProfileResponse
// @Failure 401 {object} response.ErrorBody
// @Router /api/auth/profile [get]
// Profile 返回当前登录用户的基础信息。
func (h *Handler) Profile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	response.Success(c, gin.H{
		"user_id":  userID,
		"username": username,
	})
}

// Menus 当前用户菜单
// @Summary 获取当前登录用户可见菜单
// @Tags auth
// @ID authMenus
// @Produce json
// @Security BearerAuth
// @Success 200 {object} MenusResponse
// @Failure 401 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/auth/menus [get]
// Menus 返回当前登录用户可见的菜单树。
func (h *Handler) Menus(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "用户信息无效")
		return
	}

	menus, err := h.repo.ListMenusByUserID(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询菜单失败")
		return
	}

	response.Success(c, buildMenuTree(menus))
}

// buildMenuTree 把平铺菜单组装成树形结构。
func buildMenuTree(menus []Menu) []Menu {
	menuMap := make(map[int64]*Menu, len(menus))
	for i := range menus {
		menus[i].Children = nil
		menuMap[menus[i].ID] = &menus[i]
	}

	tree := make([]Menu, 0)
	for i := range menus {
		menu := &menus[i]
		if menu.ParentID == 0 {
			tree = append(tree, *menu)
			continue
		}

		parent, ok := menuMap[menu.ParentID]
		if !ok {
			tree = append(tree, *menu)
			continue
		}

		parent.Children = append(parent.Children, *menu)
	}

	sort.Slice(tree, func(i, j int) bool {
		if tree[i].SortOrder == tree[j].SortOrder {
			return tree[i].ID < tree[j].ID
		}
		return tree[i].SortOrder < tree[j].SortOrder
	})

	for i := range tree {
		sort.Slice(tree[i].Children, func(a, b int) bool {
			if tree[i].Children[a].SortOrder == tree[i].Children[b].SortOrder {
				return tree[i].Children[a].ID < tree[i].Children[b].ID
			}
			return tree[i].Children[a].SortOrder < tree[i].Children[b].SortOrder
		})
	}

	return tree
}
