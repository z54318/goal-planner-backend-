package auth

// LoginRequest 表示登录请求体。
type LoginRequest struct {
	// Username 用户名
	Username string `json:"username"`
	// Password 密码
	Password string `json:"password"`
}

// RegisterRequest 表示注册请求体。
type RegisterRequest struct {
	// Username 用户名
	Username string `json:"username"`
	// Nickname 昵称
	Nickname string `json:"nickname"`
	// Email 邮箱
	Email string `json:"email"`
	// Password 密码
	Password string `json:"password"`
}

// User 表示登录时查询到的用户信息。
type User struct {
	// ID 用户ID
	ID int64
	// Username 用户名
	Username string
	// Nickname 昵称
	Nickname string
	// PasswordHash 密码哈希
	PasswordHash string
	// Status 用户状态
	Status string
}

// Menu 表示返回给前端的菜单节点。
type Menu struct {
	// ID 菜单ID
	ID int64 `json:"id"`
	// ParentID 父菜单ID
	ParentID int64 `json:"parent_id"`
	// Name 菜单名称
	Name string `json:"name"`
	// Path 路由路径
	Path string `json:"path"`
	// Component 组件路径
	Component string `json:"component"`
	// Icon 菜单图标
	Icon string `json:"icon"`
	// SortOrder 排序值
	SortOrder int `json:"sort_order"`
	// PermissionCode 权限码
	PermissionCode string `json:"permission_code"`
	// Hidden 是否隐藏
	Hidden bool `json:"hidden"`
	// Children 子菜单
	Children []Menu `json:"children,omitempty"`
}

// RegisterData 表示注册成功后返回的数据。
type RegisterData struct {
	// UserID 用户ID
	UserID int64 `json:"user_id"`
	// Username 用户名
	Username string `json:"username"`
	// Nickname 昵称
	Nickname string `json:"nickname"`
	// Email 邮箱
	Email string `json:"email"`
}

// RegisterResponse 表示注册成功响应。
type RegisterResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    RegisterData `json:"data"`
}

// LoginData 表示登录成功后返回的数据。
type LoginData struct {
	// Token 登录令牌
	Token string `json:"token"`
	// UserID 用户ID
	UserID int64 `json:"user_id"`
	// Username 用户名
	Username string `json:"username"`
	// Nickname 昵称
	Nickname string `json:"nickname"`
}

// LoginResponse 表示登录成功响应。
type LoginResponse struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    LoginData `json:"data"`
}

// ProfileData 表示当前登录用户基础信息。
type ProfileData struct {
	// UserID 用户ID
	UserID int64 `json:"user_id"`
	// Username 用户名
	Username string `json:"username"`
	// Nickname 昵称
	Nickname string `json:"nickname"`
}

// ProfileResponse 表示当前登录用户信息响应。
type ProfileResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    ProfileData `json:"data"`
}

// MenusResponse 表示当前用户菜单响应。
type MenusResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []Menu `json:"data"`
}
