package goal

import (
	"context"
	"database/sql"
)

// Repository 负责 goal 模块和数据库打交道。
type Repository struct {
	db *sql.DB
}

// NewRepository 创建目标仓库对象。
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create 插入一条目标记录。
func (r *Repository) Create(ctx context.Context, req CreateGoalRequest) (Goal, error) {
	query := `
		INSERT INTO goals (title, description, status)
		VALUES (?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query, req.Title, req.Description, "draft")
	if err != nil {
		return Goal{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Goal{}, err
	}

	return r.GetByID(ctx, id)
}

// GetByID 按 ID 查询单个目标。
func (r *Repository) GetByID(ctx context.Context, id int64) (Goal, error) {
	query := `
		SELECT id, title, description, status, created_at, updated_at
		FROM goals
		WHERE id = ?
	`

	var goal Goal
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&goal.ID,
		&goal.Title,
		&goal.Description,
		&goal.Status,
		&goal.CreatedAt,
		&goal.UpdatedAt,
	)
	if err != nil {
		return Goal{}, err
	}

	return goal, nil
}

// List 查询目标列表。
func (r *Repository) List(ctx context.Context) ([]Goal, error) {
	query := `
		SELECT id, title, description, status, created_at, updated_at
		FROM goals
		ORDER BY id DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	// 函数执行完毕，关闭查询结果集
	defer rows.Close()

	//创建一个长度为0的空切片
	goals := make([]Goal, 0)
	// rows.Next看结果集里面还有没有下一行，有就循环，没有就结束循环
	for rows.Next() {
		var goal Goal
		if err := rows.Scan(
			&goal.ID,
			&goal.Title,
			&goal.Description,
			&goal.Status,
			&goal.CreatedAt,
			&goal.UpdatedAt,
		); err != nil {
			return nil, err
		}

		// 将读出来的一条目标，追加到目标列表中
		goals = append(goals, goal)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return goals, nil
}
