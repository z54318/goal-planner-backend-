package plan

import "time"

// Phase 表示 plans 下的一个阶段。
type Phase struct {
	// ID 阶段ID
	ID int64 `json:"id"`
	// PlanID 所属计划ID
	PlanID int64 `json:"plan_id"`
	// Title 阶段标题
	Title string `json:"title"`
	// Description 阶段描述
	Description string `json:"description"`
	// SortOrder 阶段顺序
	SortOrder int `json:"sort_order"`
	// Tasks 任务列表
	Tasks []Task `json:"tasks,omitempty"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

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

// Task 表示阶段下的一个任务。
type Task struct {
	// ID 任务ID
	ID int64 `json:"id"`
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

// Plan 表示 plans 表中的一条计划记录。
type Plan struct {
	// ID 计划ID
	ID int64 `json:"id"`
	// UserID 用户ID
	UserID int64 `json:"user_id"`
	// GoalID 目标ID
	GoalID int64 `json:"goal_id"`
	// Title 计划标题
	Title string `json:"title"`
	// Overview 计划概述
	Overview string `json:"overview"`
	// Phases 阶段列表
	Phases []Phase `json:"phases,omitempty"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// PlanResponse 表示单个计划成功响应。
type PlanResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Plan   `json:"data"`
}

// NextStepSuggestion 表示执行建议。
type NextStepSuggestion struct {
	// Summary 执行建议摘要
	Summary string `json:"summary"`
	// NextAction 建议执行动作
	NextAction string `json:"next_action"`
	// Reason 建议原因
	Reason string `json:"reason"`
	// Checklist 执行清单
	Checklist []string `json:"checklist"`
	// Risk 主要风险
	Risk string `json:"risk"`
}

// NextStepSuggestionResponse 表示执行建议成功响应。
type NextStepSuggestionResponse struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    NextStepSuggestion `json:"data"`
}

// UpdatePlanRequest 表示更新计划请求体。
type UpdatePlanRequest struct {
	// Title 计划标题
	Title string `json:"title"`
	// Overview 计划概述
	Overview string `json:"overview"`
}

// DeletePlanRequest 表示删除计划请求体
type DeletePlanRequest struct {
	PlanID int64 `json:"plan_id"`
}
