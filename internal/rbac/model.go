package rbac

import "time"

// Role 表示一条角色记录。
type Role struct {
	// ID 角色ID
	ID int64 `json:"id"`
	// Name 角色名称
	Name string `json:"name"`
	// Code 角色编码
	Code string `json:"code"`
	// PermissionIDs 当前角色绑定的权限ID列表
	PermissionIDs []int64 `json:"permission_ids,omitempty"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateRoleRequest 表示新增角色请求体。
type CreateRoleRequest struct {
	// Name 角色名称
	Name string `json:"name"`
	// Code 角色编码
	Code string `json:"code"`
}

// Permission 表示一条权限记录。
type Permission struct {
	// ID 权限ID
	ID int64 `json:"id"`
	// Name 权限名称
	Name string `json:"name"`
	// Code 权限编码
	Code string `json:"code"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// CreatePermissionRequest 表示新增权限请求体。
type CreatePermissionRequest struct {
	// Name 权限名称
	Name string `json:"name"`
	// Code 权限编码
	Code string `json:"code"`
}

// UpdatePermissionRequest 表示更新权限请求体。
type UpdatePermissionRequest struct {
	// Name 权限名称
	Name string `json:"name"`
	// Code 权限编码
	Code string `json:"code"`
}

// UpdateRolePermissionsRequest 表示更新角色权限绑定请求体。
type UpdateRolePermissionsRequest struct {
	// PermissionIDs 权限ID列表
	PermissionIDs []int64 `json:"permission_ids"`
}

// RoleListResponse 表示角色列表成功响应。
type RoleListResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []Role `json:"data"`
}

// RoleResponse 表示单个角色成功响应。
type RoleResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Role   `json:"data"`
}

// PermissionListResponse 表示权限列表成功响应。
type PermissionListResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    []Permission `json:"data"`
}

// PermissionResponse 表示单个权限成功响应。
type PermissionResponse struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Data    Permission `json:"data"`
}

// RolePermissionIDsData 表示角色权限ID响应数据。
type RolePermissionIDsData struct {
	// RoleID 角色ID
	RoleID int64 `json:"role_id"`
	// PermissionIDs 权限ID列表
	PermissionIDs []int64 `json:"permission_ids"`
}

// RolePermissionIDsResponse 表示角色权限ID成功响应。
type RolePermissionIDsResponse struct {
	Code    int                   `json:"code"`
	Message string                `json:"message"`
	Data    RolePermissionIDsData `json:"data"`
}
