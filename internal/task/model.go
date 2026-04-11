package task

import "time"

// TaskStatus 表示任务状态。
type TaskStatus string

const (
	// TaskStatusTodo 待开始
	TaskStatusTodo TaskStatus = "todo"
	// TaskStatusInProgress 进行中
	TaskStatusInProgress TaskStatus = "in_progress"
	// TaskStatusDone 已完成
	TaskStatusDone TaskStatus = "done"
)

// TaskPriority 表示任务优先级。
type TaskPriority string

const (
	// TaskPriorityHigh 高优先级
	TaskPriorityHigh TaskPriority = "high"
	// TaskPriorityMedium 中优先级
	TaskPriorityMedium TaskPriority = "medium"
	// TaskPriorityLow 低优先级
	TaskPriorityLow TaskPriority = "low"
)

// Task 表示一条可执行任务。
type Task struct {
	// GoalID 所属目标ID
	GoalID int64 `json:"goal_id,omitempty"`
	// PlanID 所属计划ID
	PlanID int64 `json:"plan_id,omitempty"`
	// ID 任务ID
	ID int64 `json:"id"`
	// PhaseID 所属阶段ID
	PhaseID int64 `json:"phase_id"`
	// GoalTitle 目标标题
	GoalTitle string `json:"goal_title,omitempty"`
	// PhaseTitle 阶段标题
	PhaseTitle string `json:"phase_title,omitempty"`
	// Title 任务标题
	Title string `json:"title"`
	// Description 任务描述
	Description string `json:"description"`
	// EstimatedDays 预估天数
	EstimatedDays int `json:"estimated_days"`
	// Deliverables 交付物
	Deliverables string `json:"deliverables"`
	// Deadline 任务截止时间
	Deadline *time.Time `json:"deadline,omitempty"`
	// Priority 任务优先级
	Priority TaskPriority `json:"priority" enums:"high,medium,low"`
	// Status 任务状态
	Status TaskStatus `json:"status" enums:"todo,in_progress,done"`
	// SortOrder 任务顺序
	SortOrder int `json:"sort_order"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// TaskResponse 表示单个任务成功响应。
type TaskResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Task   `json:"data"`
}

// TaskListData 表示任务列表响应数据。
type TaskListData struct {
	// List 任务列表
	List []Task `json:"list"`
	// Total 总数量
	Total int `json:"total"`
	// Page 当前页
	Page int `json:"page"`
	// PageSize 每页条数
	PageSize int `json:"page_size"`
}

// TaskListResponse 表示任务列表成功响应。
type TaskListResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    TaskListData `json:"data"`
}

// CreateTaskRequest 表示新增任务请求体。
type CreateTaskRequest struct {
	// PhaseID 所属阶段ID
	PhaseID int64 `json:"phase_id"`
	// Title 任务标题
	Title string `json:"title"`
	// Description 任务描述
	Description string `json:"description"`
	// EstimatedDays 预估天数
	EstimatedDays int `json:"estimated_days"`
	// Deliverables 交付物
	Deliverables string `json:"deliverables"`
	// Deadline 任务截止时间，使用 RFC3339 格式
	Deadline *time.Time `json:"deadline"`
	// Priority 任务优先级：high 高，medium 中，low 低
	Priority TaskPriority `json:"priority" enums:"high,medium,low"`
	// SortOrder 任务顺序
	SortOrder int `json:"sort_order"`
}

// UpdateTaskRequest 表示编辑任务请求体。
type UpdateTaskRequest struct {
	// PhaseID 所属阶段ID
	PhaseID int64 `json:"phase_id"`
	// Title 任务标题
	Title string `json:"title"`
	// Description 任务描述
	Description string `json:"description"`
	// EstimatedDays 预估天数
	EstimatedDays int `json:"estimated_days"`
	// Deliverables 交付物
	Deliverables string `json:"deliverables"`
	// Deadline 任务截止时间，使用 RFC3339 格式
	Deadline *time.Time `json:"deadline"`
	// Priority 任务优先级：high 高，medium 中，low 低
	Priority TaskPriority `json:"priority" enums:"high,medium,low"`
	// SortOrder 任务顺序
	SortOrder int `json:"sort_order"`
}

// ListTasksRequest 表示任务列表查询参数。
type ListTasksRequest struct {
	// Status 任务状态筛选
	Status TaskStatus `form:"status" json:"status"`
	// GoalID 目标ID筛选
	GoalID int64 `form:"goal_id" json:"goal_id"`
	// PhaseID 阶段ID筛选
	PhaseID int64 `form:"phase_id" json:"phase_id"`
	// Page 页码，从1开始
	Page int `form:"page" json:"page"`
	// PageSize 每页条数
	PageSize int `form:"page_size" json:"page_size"`
}

// UpdateTaskStatusRequest 表示更新任务状态请求体。
type UpdateTaskStatusRequest struct {
	// Status 任务状态：todo 待开始，in_progress 进行中，done 已完成
	Status TaskStatus `json:"status" enums:"todo,in_progress,done"`
}

// DeleteTaskResponse 表示删除任务成功响应。
type DeleteTaskResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    DeleteTaskData `json:"data"`
}

// DeleteTaskData 表示删除任务响应数据。
type DeleteTaskData struct {
	// Deleted 是否已删除
	Deleted bool `json:"deleted"`
}
