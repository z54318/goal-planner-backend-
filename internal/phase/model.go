package phase

import "time"

// Phase 表示一条计划阶段。
type Phase struct {
	// GoalID 所属目标ID
	GoalID int64 `json:"goal_id,omitempty"`
	// PlanID 所属计划ID
	PlanID int64 `json:"plan_id"`
	// ID 阶段ID
	ID int64 `json:"id"`
	// Title 阶段标题
	Title string `json:"title"`
	// Description 阶段描述
	Description string `json:"description"`
	// SortOrder 阶段顺序
	SortOrder int `json:"sort_order"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// PhaseResponse 表示单个阶段成功响应。
type PhaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Phase  `json:"data"`
}

// UpdatePhaseRequest 表示编辑阶段请求体。
type UpdatePhaseRequest struct {
	// Title 阶段标题
	Title string `json:"title"`
	// Description 阶段描述
	Description string `json:"description"`
	// SortOrder 阶段顺序
	SortOrder int `json:"sort_order"`
}
