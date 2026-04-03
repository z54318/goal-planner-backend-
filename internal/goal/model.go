package goal

import "time"

// Goal 表示 goals 表中的一条目标记录。
type Goal struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateGoalRequest 表示创建目标时前端传来的请求体。
type CreateGoalRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// GoalResponse 表示单个目标成功响应。
type GoalResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Goal   `json:"data"`
}

// GoalListData 表示目标列表响应中的 data 字段。
type GoalListData struct {
	List  []Goal `json:"list"`
	Total int    `json:"total"`
}

// GoalListResponse 表示目标列表成功响应。
type GoalListResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    GoalListData `json:"data"`
}
