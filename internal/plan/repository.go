package plan

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"time"

	appai "goal-planner/internal/infra/ai"
)

// Repository 负责 plan 模块和数据库打交道。
type Repository struct {
	db *sql.DB
}

// NewRepository 创建计划仓库对象。
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// GetGoalForGeneration 获取生成计划所需的目标信息。
func (r *Repository) GetGoalForGeneration(ctx context.Context, userID int64, goalID int64) (appai.GoalInput, error) {
	query := `
		SELECT title, description, category, target_deadline
		FROM goals
		WHERE id = ? AND user_id = ?
	`

	var goal appai.GoalInput
	err := r.db.QueryRowContext(ctx, query, goalID, userID).Scan(
		&goal.Title,
		&goal.Description,
		&goal.Category,
		&goal.TargetDeadline,
	)
	if err != nil {
		return appai.GoalInput{}, err
	}

	return goal, nil
}

// GetByGoalID 按目标ID查询当前用户的一条计划记录。
func (r *Repository) GetByGoalID(ctx context.Context, userID int64, goalID int64) (Plan, error) {
	query := `
		SELECT id, user_id, goal_id, title, overview, created_at, updated_at
		FROM plans
		WHERE goal_id = ? AND user_id = ?
	`

	var plan Plan
	err := r.db.QueryRowContext(ctx, query, goalID, userID).Scan(
		&plan.ID,
		&plan.UserID,
		&plan.GoalID,
		&plan.Title,
		&plan.Overview,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)
	if err != nil {
		return Plan{}, err
	}

	phases, err := r.listPhasesByPlanID(ctx, plan.ID)
	if err != nil {
		return Plan{}, err
	}
	plan.Phases = phases

	return plan, nil
}

// UpdateByGoalID 更新当前用户某个目标下的计划基础信息。
func (r *Repository) UpdateByGoalID(ctx context.Context, userID int64, goalID int64, req UpdatePlanRequest) (Plan, error) {
	result, err := r.db.ExecContext(
		ctx,
		`
			UPDATE plans
			SET title = ?, overview = ?
			WHERE goal_id = ? AND user_id = ?
		`,
		req.Title,
		req.Overview,
		goalID,
		userID,
	)
	if err != nil {
		return Plan{}, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return Plan{}, err
	}
	if affected == 0 {
		return Plan{}, sql.ErrNoRows
	}

	return r.GetByGoalID(ctx, userID, goalID)
}

// DeleteByGoalID 删除当前用户某个目标下的计划及其关联数据。
func (r *Repository) DeleteByGoalID(ctx context.Context, userID int64, goalID int64) error {
	// 开启一个数据库事务：将数据库操作绑一起，要么全成功，要么全失败
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// 回滚事务：defer等当前函数结束时再执行
	defer tx.Rollback()

	var planID int64
	err = tx.QueryRowContext(
		ctx,
		`SELECT id FROM plans WHERE goal_id = ? AND user_id = ?`,
		goalID,
		userID,
	).Scan(&planID)
	// Scan:将数据库查出来的列值读出来填到变量里 &：表示取这个变量的地址，让它能往里面写值 如果不写&就是变量里的值
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`DELETE FROM tasks WHERE phase_id IN (SELECT id FROM phases WHERE plan_id = ? )`,
		planID,
	)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(
		ctx,
		`DELETE FROM phases WHERE plan_id = ?`,
		planID,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`DELETE FROM plans WHERE id = ? AND user_id = ?`,
		planID,
		userID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// CreateGenerated 为当前用户的目标创建一条 AI 生成的计划记录。
func (r *Repository) CreateGenerated(ctx context.Context, userID int64, goalID int64, output appai.PlanOutput) (Plan, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Plan{}, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO plans (user_id, goal_id, title, overview)
		SELECT ?, id, ?, ?
		FROM goals
		WHERE id = ? AND user_id = ?
	`

	result, err := tx.ExecContext(ctx, query, userID, output.Title, output.Overview, goalID, userID)
	if err != nil {
		return Plan{}, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return Plan{}, err
	}
	if affected == 0 {
		return Plan{}, sql.ErrNoRows
	}

	planID, err := result.LastInsertId()
	if err != nil {
		return Plan{}, err
	}

	for _, phase := range output.Phases {
		phaseResult, err := tx.ExecContext(
			ctx,
			`INSERT INTO phases (plan_id, title, description, sort_order) VALUES (?, ?, ?, ?)`,
			planID,
			phase.Title,
			phase.Description,
			phase.Order,
		)
		if err != nil {
			return Plan{}, err
		}

		phaseID, err := phaseResult.LastInsertId()
		if err != nil {
			return Plan{}, err
		}

		for _, task := range phase.Tasks {
			deadline, err := parseGeneratedTaskDeadline(task.Deadline)
			if err != nil {
				return Plan{}, err
			}

			if _, err := tx.ExecContext(
				ctx,
				`INSERT INTO tasks (phase_id, title, description, estimated_days, deliverables, deadline, priority, status, sort_order) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				phaseID,
				task.Title,
				task.Description,
				task.EstimatedDays,
				task.Deliverables,
				deadline,
				task.Priority,
				"todo",
				task.Order,
			); err != nil {
				return Plan{}, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return Plan{}, err
	}

	return r.GetByGoalID(ctx, userID, goalID)
}

// RegenerateGenerated 删除旧计划后重新生成一条计划记录。
func (r *Repository) RegenerateGenerated(ctx context.Context, userID int64, goalID int64, output appai.PlanOutput) (Plan, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Plan{}, err
	}
	defer tx.Rollback()

	var existingPlanID int64
	err = tx.QueryRowContext(
		ctx,
		`SELECT id FROM plans WHERE goal_id = ? AND user_id = ?`,
		goalID,
		userID,
	).Scan(&existingPlanID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return Plan{}, err
	}

	if err == nil {
		if _, err := tx.ExecContext(ctx, `DELETE FROM tasks WHERE phase_id IN (SELECT id FROM phases WHERE plan_id = ?)`, existingPlanID); err != nil {
			return Plan{}, err
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM phases WHERE plan_id = ?`, existingPlanID); err != nil {
			return Plan{}, err
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM plans WHERE id = ?`, existingPlanID); err != nil {
			return Plan{}, err
		}
	}

	query := `
		INSERT INTO plans (user_id, goal_id, title, overview)
		SELECT ?, id, ?, ?
		FROM goals
		WHERE id = ? AND user_id = ?
	`

	result, err := tx.ExecContext(ctx, query, userID, output.Title, output.Overview, goalID, userID)
	if err != nil {
		return Plan{}, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return Plan{}, err
	}
	if affected == 0 {
		return Plan{}, sql.ErrNoRows
	}

	planID, err := result.LastInsertId()
	if err != nil {
		return Plan{}, err
	}

	for _, phase := range output.Phases {
		phaseResult, err := tx.ExecContext(
			ctx,
			`INSERT INTO phases (plan_id, title, description, sort_order) VALUES (?, ?, ?, ?)`,
			planID,
			phase.Title,
			phase.Description,
			phase.Order,
		)
		if err != nil {
			return Plan{}, err
		}

		phaseID, err := phaseResult.LastInsertId()
		if err != nil {
			return Plan{}, err
		}

		for _, task := range phase.Tasks {
			deadline, err := parseGeneratedTaskDeadline(task.Deadline)
			if err != nil {
				return Plan{}, err
			}

			if _, err := tx.ExecContext(
				ctx,
				`INSERT INTO tasks (phase_id, title, description, estimated_days, deliverables, deadline, priority, status, sort_order) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				phaseID,
				task.Title,
				task.Description,
				task.EstimatedDays,
				task.Deliverables,
				deadline,
				task.Priority,
				"todo",
				task.Order,
			); err != nil {
				return Plan{}, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return Plan{}, err
	}

	return r.GetByGoalID(ctx, userID, goalID)
}

func (r *Repository) listPhasesByPlanID(ctx context.Context, planID int64) ([]Phase, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, plan_id, title, description, sort_order, created_at, updated_at FROM phases WHERE plan_id = ? ORDER BY sort_order ASC, id ASC`,
		planID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	phases := make([]Phase, 0)
	phaseIDs := make([]int64, 0)
	for rows.Next() {
		var phase Phase
		if err := rows.Scan(
			&phase.ID,
			&phase.PlanID,
			&phase.Title,
			&phase.Description,
			&phase.SortOrder,
			&phase.CreatedAt,
			&phase.UpdatedAt,
		); err != nil {
			return nil, err
		}

		phases = append(phases, phase)
		phaseIDs = append(phaseIDs, phase.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(phases) == 0 {
		return phases, nil
	}

	taskMap, err := r.listTasksByPhaseIDs(ctx, phaseIDs)
	if err != nil {
		return nil, err
	}

	for i := range phases {
		phases[i].Tasks = taskMap[phases[i].ID]
	}

	return phases, nil
}

func (r *Repository) listTasksByPhaseIDs(ctx context.Context, phaseIDs []int64) (map[int64][]Task, error) {
	query := `
		SELECT id, phase_id, title, description, estimated_days, deliverables, deadline, priority, status, sort_order, created_at, updated_at
		FROM tasks
		WHERE phase_id = ?
		ORDER BY sort_order ASC, id ASC
	`

	taskMap := make(map[int64][]Task, len(phaseIDs))
	for _, phaseID := range phaseIDs {
		rows, err := r.db.QueryContext(ctx, query, phaseID)
		if err != nil {
			return nil, err
		}

		tasks := make([]Task, 0)
		for rows.Next() {
			var task Task
			if err := rows.Scan(
				&task.ID,
				&task.PhaseID,
				&task.Title,
				&task.Description,
				&task.EstimatedDays,
				&task.Deliverables,
				&task.Deadline,
				&task.Priority,
				&task.Status,
				&task.SortOrder,
				&task.CreatedAt,
				&task.UpdatedAt,
			); err != nil {
				rows.Close()
				return nil, err
			}

			tasks = append(tasks, task)
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return nil, err
		}
		rows.Close()
		taskMap[phaseID] = tasks
	}

	for phaseID := range taskMap {
		sort.Slice(taskMap[phaseID], func(i, j int) bool {
			if taskMap[phaseID][i].SortOrder == taskMap[phaseID][j].SortOrder {
				return taskMap[phaseID][i].ID < taskMap[phaseID][j].ID
			}
			return taskMap[phaseID][i].SortOrder < taskMap[phaseID][j].SortOrder
		})
	}

	return taskMap, nil
}

func parseGeneratedTaskDeadline(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}

	deadline, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, err
	}

	return &deadline, nil
}
