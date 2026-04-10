package menu

import "time"

// Menu 表示后台维护的一条菜单记录。
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
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateMenuRequest 表示新增菜单请求体。
type CreateMenuRequest struct {
	// ParentID 父菜单ID，0 表示顶级
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
}

// UpdateMenuRequest 表示更新菜单请求体。
type UpdateMenuRequest struct {
	// ParentID 父菜单ID，0 表示顶级
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
}

// MenuListResponse 表示菜单列表成功响应。
type MenuListResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []Menu `json:"data"`
}

// MenuResponse 表示单个菜单成功响应。
type MenuResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Menu   `json:"data"`
}
