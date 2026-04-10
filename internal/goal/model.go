package goal

import "time"

// GoalStatus 目标状态。
// - draft: 未执行
// - active: 执行中
// - completed: 已完成
// - archived: 已归档
type GoalStatus string

const (
	// GoalStatusDraft 未执行
	GoalStatusDraft GoalStatus = "draft"
	// GoalStatusActive 执行中
	GoalStatusActive GoalStatus = "active"
	// GoalStatusCompleted 已完成
	GoalStatusCompleted GoalStatus = "completed"
	// GoalStatusArchived 已归档
	GoalStatusArchived GoalStatus = "archived"
)

// Goal 表示 goals 表中的一条目标记录。
type Goal struct {
	// ID 目标ID
	ID int64 `json:"id"`
	// UserID 用户ID
	UserID int64 `json:"user_id"`
	// Title 目标标题
	Title string `json:"title"`
	// Description 目标描述
	Description string `json:"description"`
	// Category 目标分类
	Category string `json:"category"`
	// TargetDeadline 截止时间
	TargetDeadline *time.Time `json:"target_deadline,omitempty"`
	// Status 目标状态
	Status GoalStatus `json:"status" enums:"draft,active,completed,archived"`
	// AggregateStatus 按任务聚合出来的目标状态
	AggregateStatus GoalStatus `json:"aggregate_status" enums:"draft,active,completed,archived"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateGoalRequest 表示创建目标时前端传来的请求体。
type CreateGoalRequest struct {
	// Title 目标标题
	Title string `json:"title"`
	// Description 目标描述
	Description string `json:"description"`
	// Category 目标分类
	Category string `json:"category"`
	// TargetDeadline 截止时间
	TargetDeadline *time.Time `json:"target_deadline"`
}

// UpdateGoalRequest 表示更新目标时前端传来的请求体。
type UpdateGoalRequest struct {
	// Title 目标标题
	Title string `json:"title"`
	// Description 目标描述
	Description string `json:"description"`
	// Category 目标分类
	Category string `json:"category"`
	// TargetDeadline 截止时间
	TargetDeadline *time.Time `json:"target_deadline"`
}

// UpdateGoalStatusRequest 表示更新目标状态时前端传来的请求体。
type UpdateGoalStatusRequest struct {
	// Status 目标状态
	Status GoalStatus `json:"status" enums:"draft,active,completed,archived"`
}

// GoalResponse 表示单个目标成功响应。
type GoalResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Goal   `json:"data"`
}

// ListGoalsRequest 表示目标列表查询参数。
type ListGoalsRequest struct {
	// Page 页码，从1开始
	Page int `form:"page" json:"page"`
	// PageSize 每页条数
	PageSize int `form:"page_size" json:"page_size"`
}

// GoalListData 表示目标列表响应中的 data 字段。
type GoalListData struct {
	// List 目标列表
	List []Goal `json:"list"`
	// Total 总数
	Total int `json:"total"`
	// Page 当前页
	Page int `json:"page"`
	// PageSize 每页条数
	PageSize int `json:"page_size"`
}

// GoalListResponse 表示目标列表成功响应。
type GoalListResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    GoalListData `json:"data"`
}
