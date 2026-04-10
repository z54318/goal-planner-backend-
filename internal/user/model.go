package user

import "time"

// Role 表示用户已绑定的一条角色摘要。
type Role struct {
	// ID 角色ID
	ID int64 `json:"id"`
	// Name 角色名称
	Name string `json:"name"`
	// Code 角色编码
	Code string `json:"code"`
}

// User 表示后台管理中的一条用户记录。
type User struct {
	// ID 用户ID
	ID int64 `json:"id"`
	// Username 用户名
	Username string `json:"username"`
	// Nickname 昵称
	Nickname string `json:"nickname"`
	// Email 邮箱
	Email string `json:"email"`
	// Status 用户状态
	Status string `json:"status"`
	// RoleIDs 已绑定角色ID
	RoleIDs []int64 `json:"role_ids,omitempty"`
	// Roles 已绑定角色列表
	Roles []Role `json:"roles,omitempty"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// UserListResponse 表示用户列表成功响应。
type UserListResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []User `json:"data"`
}

// UserResponse 表示单个用户成功响应。
type UserResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    User   `json:"data"`
}

// UpdateUserRolesRequest 表示更新用户角色绑定请求体。
type UpdateUserRolesRequest struct {
	// RoleIDs 角色ID列表
	RoleIDs []int64 `json:"role_ids"`
}
