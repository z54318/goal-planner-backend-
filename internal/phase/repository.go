package phase

import (
	"context"
	"database/sql"
)

// Repository 负责 phase 模块和数据库打交道。
type Repository struct {
	db *sql.DB
}

// NewRepository 创建阶段仓库对象。
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// GetByID 按阶段ID查询当前用户的阶段详情。
func (r *Repository) GetByID(ctx context.Context, userID int64, phaseID int64) (Phase, error) {
	query := `
		SELECT
			pl.goal_id,
			p.plan_id,
			p.id,
			p.title,
			p.description,
			p.sort_order,
			p.created_at,
			p.updated_at
		FROM phases p
		INNER JOIN plans pl ON pl.id = p.plan_id
		WHERE p.id = ? AND pl.user_id = ?
	`

	var phase Phase
	err := r.db.QueryRowContext(ctx, query, phaseID, userID).Scan(
		&phase.GoalID,
		&phase.PlanID,
		&phase.ID,
		&phase.Title,
		&phase.Description,
		&phase.SortOrder,
		&phase.CreatedAt,
		&phase.UpdatedAt,
	)
	if err != nil {
		return Phase{}, err
	}

	return phase, nil
}

// Update 更新当前用户的一条阶段。
func (r *Repository) Update(ctx context.Context, userID int64, phaseID int64, req UpdatePhaseRequest) (Phase, error) {
	result, err := r.db.ExecContext(
		ctx,
		`
			UPDATE phases p
			INNER JOIN plans pl ON pl.id = p.plan_id
			SET
				p.title = ?,
				p.description = ?,
				p.sort_order = ?
			WHERE p.id = ? AND pl.user_id = ?
		`,
		req.Title,
		req.Description,
		req.SortOrder,
		phaseID,
		userID,
	)
	if err != nil {
		return Phase{}, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return Phase{}, err
	}
	if affected == 0 {
		return Phase{}, sql.ErrNoRows
	}

	return r.GetByID(ctx, userID, phaseID)
}
