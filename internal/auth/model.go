package auth

// LoginRequest 表示登录请求体。
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterRequest 表示注册请求体。
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User 表示登录时查询到的用户信息。
type User struct {
	ID           int64
	Username     string
	PasswordHash string
	Status       string
}

// Menu 表示返回给前端的菜单节点。
type Menu struct {
	ID             int64  `json:"id"`
	ParentID       int64  `json:"parent_id"`
	Name           string `json:"name"`
	Path           string `json:"path"`
	Component      string `json:"component"`
	Icon           string `json:"icon"`
	SortOrder      int    `json:"sort_order"`
	PermissionCode string `json:"permission_code"`
	Hidden         bool   `json:"hidden"`
	Children       []Menu `json:"children,omitempty"`
}

// RegisterData 表示注册成功后返回的数据。
type RegisterData struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// RegisterResponse 表示注册成功响应。
type RegisterResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    RegisterData `json:"data"`
}

// LoginData 表示登录成功后返回的数据。
type LoginData struct {
	Token    string `json:"token"`
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

// LoginResponse 表示登录成功响应。
type LoginResponse struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    LoginData `json:"data"`
}

// ProfileData 表示当前登录用户基础信息。
type ProfileData struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
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
