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
func (r *Repository) Create(ctx context.Context, userID int64, req CreateGoalRequest) (Goal, error) {
	query := `
		INSERT INTO goals (user_id, title, description, category, target_deadline, status)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		userID,
		req.Title,
		req.Description,
		req.Category,
		req.TargetDeadline,
		"draft",
	)
	if err != nil {
		return Goal{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Goal{}, err
	}

	return r.GetByID(ctx, userID, id)
}

// GetByID 按 ID 查询当前用户的一条目标记录。
func (r *Repository) GetByID(ctx context.Context, userID int64, id int64) (Goal, error) {
	query := `
		SELECT
			g.id,
			g.user_id,
			g.title,
			g.description,
			g.category,
			g.target_deadline,
			g.status,
			CASE
				WHEN COUNT(t.id) = 0 THEN 'draft'
				WHEN SUM(CASE WHEN t.status = 'in_progress' THEN 1 ELSE 0 END) > 0 THEN 'active'
				WHEN SUM(CASE WHEN t.status = 'done' THEN 1 ELSE 0 END) = COUNT(t.id) THEN 'completed'
				WHEN SUM(CASE WHEN t.status = 'todo' THEN 1 ELSE 0 END) = COUNT(t.id) THEN 'draft'
				ELSE 'active'
			END AS aggregate_status,
			g.created_at,
			g.updated_at
		FROM goals g
		LEFT JOIN plans p ON p.goal_id = g.id AND p.user_id = g.user_id
		LEFT JOIN phases ph ON ph.plan_id = p.id
		LEFT JOIN tasks t ON t.phase_id = ph.id
		WHERE g.id = ? AND g.user_id = ?
		GROUP BY
			g.id,
			g.user_id,
			g.title,
			g.description,
			g.category,
			g.target_deadline,
			g.status,
			g.created_at,
			g.updated_at
	`

	var goal Goal
	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(
		&goal.ID,
		&goal.UserID,
		&goal.Title,
		&goal.Description,
		&goal.Category,
		&goal.TargetDeadline,
		&goal.Status,
		&goal.AggregateStatus,
		&goal.CreatedAt,
		&goal.UpdatedAt,
	)
	if err != nil {
		return Goal{}, err
	}

	return goal, nil
}

// ListByUserID 查询当前用户的目标列表。
func (r *Repository) ListByUserID(ctx context.Context, userID int64, req ListGoalsRequest) ([]Goal, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM goals WHERE user_id = ?`, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT
			g.id,
			g.user_id,
			g.title,
			g.description,
			g.category,
			g.target_deadline,
			g.status,
			CASE
				WHEN COUNT(t.id) = 0 THEN 'draft'
				WHEN SUM(CASE WHEN t.status = 'in_progress' THEN 1 ELSE 0 END) > 0 THEN 'active'
				WHEN SUM(CASE WHEN t.status = 'done' THEN 1 ELSE 0 END) = COUNT(t.id) THEN 'completed'
				WHEN SUM(CASE WHEN t.status = 'todo' THEN 1 ELSE 0 END) = COUNT(t.id) THEN 'draft'
				ELSE 'active'
			END AS aggregate_status,
			g.created_at,
			g.updated_at
		FROM goals g
		LEFT JOIN plans p ON p.goal_id = g.id AND p.user_id = g.user_id
		LEFT JOIN phases ph ON ph.plan_id = p.id
		LEFT JOIN tasks t ON t.phase_id = ph.id
		WHERE g.user_id = ?
		GROUP BY
			g.id,
			g.user_id,
			g.title,
			g.description,
			g.category,
			g.target_deadline,
			g.status,
			g.created_at,
			g.updated_at
		ORDER BY g.id DESC
		LIMIT ? OFFSET ?
	`

	offset := (req.Page - 1) * req.PageSize
	rows, err := r.db.QueryContext(ctx, query, userID, req.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	goals := make([]Goal, 0)
	for rows.Next() {
		var goal Goal
		if err := rows.Scan(
			&goal.ID,
			&goal.UserID,
			&goal.Title,
			&goal.Description,
			&goal.Category,
			&goal.TargetDeadline,
			&goal.Status,
			&goal.AggregateStatus,
			&goal.CreatedAt,
			&goal.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		goals = append(goals, goal)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return goals, total, nil
}

// Update 更新当前用户的一条目标记录。
func (r *Repository) Update(ctx context.Context, userID int64, id int64, req UpdateGoalRequest) (Goal, error) {
	query := `
		UPDATE goals
		SET title = ?, description = ?, category = ?, target_deadline = ?
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		req.Title,
		req.Description,
		req.Category,
		req.TargetDeadline,
		id,
		userID,
	)
	if err != nil {
		return Goal{}, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return Goal{}, err
	}
	if affected == 0 {
		return Goal{}, sql.ErrNoRows
	}

	return r.GetByID(ctx, userID, id)
}

// UpdateStatus 更新当前用户目标的状态。
func (r *Repository) UpdateStatus(ctx context.Context, userID int64, id int64, status GoalStatus) (Goal, error) {
	query := `
		UPDATE goals
		SET status = ?
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.ExecContext(ctx, query, status, id, userID)
	if err != nil {
		return Goal{}, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return Goal{}, err
	}
	if affected == 0 {
		return Goal{}, sql.ErrNoRows
	}

	return r.GetByID(ctx, userID, id)
}

// Delete 删除当前用户的一条目标记录。
func (r *Repository) Delete(ctx context.Context, userID int64, id int64) error {
	query := `
		DELETE FROM goals
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
