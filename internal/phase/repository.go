package phase

import (
	"context"
	"database/sql"
	"time"

	storedsuggestion "goal-planner/internal/common/suggestion"
	appai "goal-planner/internal/infra/ai"
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

// GetSuggestionContextByID 查询阶段下一步建议所需上下文。
func (r *Repository) GetSuggestionContextByID(ctx context.Context, userID int64, phaseID int64) (appai.PhaseSuggestionInput, error) {
	query := `
		SELECT
			g.title,
			pl.title,
			p.title,
			p.description
		FROM phases p
		INNER JOIN plans pl ON pl.id = p.plan_id
		INNER JOIN goals g ON g.id = pl.goal_id
		WHERE p.id = ? AND pl.user_id = ?
	`

	var input appai.PhaseSuggestionInput
	err := r.db.QueryRowContext(ctx, query, phaseID, userID).Scan(
		&input.GoalTitle,
		&input.PlanTitle,
		&input.PhaseTitle,
		&input.PhaseDescription,
	)
	if err != nil {
		return appai.PhaseSuggestionInput{}, err
	}

	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT title, description, status, priority, deadline
			FROM tasks
			WHERE phase_id = ?
			ORDER BY sort_order ASC, id ASC
		`,
		phaseID,
	)
	if err != nil {
		return appai.PhaseSuggestionInput{}, err
	}
	defer rows.Close()

	input.Tasks = make([]appai.SuggestionTaskDigest, 0)
	for rows.Next() {
		var task appai.SuggestionTaskDigest
		var deadline sql.NullTime
		if err := rows.Scan(
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&deadline,
		); err != nil {
			return appai.PhaseSuggestionInput{}, err
		}
		task.Deadline = formatNullDeadline(deadline)
		input.Tasks = append(input.Tasks, task)
	}
	if err := rows.Err(); err != nil {
		return appai.PhaseSuggestionInput{}, err
	}

	return input, nil
}

// SaveSuggestionByID 保存一条阶段执行建议。
func (r *Repository) SaveSuggestionByID(ctx context.Context, userID int64, phaseID int64, item appai.NextStepSuggestion) error {
	if _, err := r.getOwnedPhaseID(ctx, userID, phaseID); err != nil {
		return err
	}

	return storedsuggestion.Upsert(ctx, r.db, userID, storedsuggestion.TargetTypePhase, phaseID, item)
}

// GetSavedSuggestionByID 查询一条已保存的阶段执行建议。
func (r *Repository) GetSavedSuggestionByID(ctx context.Context, userID int64, phaseID int64) (appai.NextStepSuggestion, error) {
	if _, err := r.getOwnedPhaseID(ctx, userID, phaseID); err != nil {
		return appai.NextStepSuggestion{}, err
	}

	return storedsuggestion.Get(ctx, r.db, userID, storedsuggestion.TargetTypePhase, phaseID)
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

// Delete 删除当前用户的一条阶段
func (r *Repository) Delete(ctx context.Context, userID int64, phaseID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var ownedPhaseID int64
	err = tx.QueryRowContext(
		ctx,
		`SELECT p.id FROM phases p INNER JOIN plans pl ON pl.id = p.plan_id WHERE p.id = ? AND pl.user_id = ?`,
		phaseID,
		userID,
	).Scan(&ownedPhaseID)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(
		ctx,
		`
			DELETE FROM ai_suggestions
			WHERE user_id = ? AND target_type = ? AND target_id IN (
				SELECT id FROM tasks WHERE phase_id = ?
			)
		`,
		userID,
		string(storedsuggestion.TargetTypeTask),
		ownedPhaseID,
	); err != nil {
		return err
	}

	if _, err := tx.ExecContext(
		ctx,
		`
			DELETE FROM ai_suggestions
			WHERE user_id = ? AND target_type = ? AND target_id = ?
		`,
		userID,
		string(storedsuggestion.TargetTypePhase),
		ownedPhaseID,
	); err != nil {
		return err
	}

	if _, err := tx.ExecContext(
		ctx,
		`DELETE FROM tasks WHERE phase_id = ? `,
		ownedPhaseID,
	); err != nil {
		return err
	}

	result, err := tx.ExecContext(
		ctx,
		`DELETE FROM phases WHERE id = ?`,
		ownedPhaseID,
	)
	if err != nil {
		return err
	}

	// RowsAffected:拿到result这个sql语句影响了多少行
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	// 删除成功后正式提交业务
	return tx.Commit()

}

func (r *Repository) getOwnedPhaseID(ctx context.Context, userID int64, phaseID int64) (int64, error) {
	var ownedPhaseID int64
	err := r.db.QueryRowContext(
		ctx,
		`SELECT p.id FROM phases p INNER JOIN plans pl ON pl.id = p.plan_id WHERE p.id = ? AND pl.user_id = ?`,
		phaseID,
		userID,
	).Scan(&ownedPhaseID)
	if err != nil {
		return 0, err
	}

	return ownedPhaseID, nil
}

func formatNullDeadline(deadline sql.NullTime) string {
	if !deadline.Valid {
		return ""
	}
	return deadline.Time.Format(time.RFC3339)
}
