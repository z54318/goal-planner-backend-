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
