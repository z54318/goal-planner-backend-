package task

import (
	"context"
	"database/sql"
	"strings"
	"time"

	storedsuggestion "goal-planner/internal/common/suggestion"
	appai "goal-planner/internal/infra/ai"
)

// Repository 负责 task 模块和数据库打交道。
type Repository struct {
	db *sql.DB
}

// NewRepository 创建任务仓库对象。
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// GetByID 按任务ID查询当前用户的任务详情。
func (r *Repository) GetByID(ctx context.Context, userID int64, taskID int64) (Task, error) {
	query := `
		SELECT
			pl.goal_id,
			p.plan_id,
			t.id,
			t.phase_id,
			g.title,
			p.title,
			t.title,
			t.description,
			t.estimated_days,
			t.deliverables,
			t.deadline,
			t.priority,
			t.status,
			t.sort_order,
			t.created_at,
			t.updated_at
		FROM tasks t
		INNER JOIN phases p ON p.id = t.phase_id
		INNER JOIN plans pl ON pl.id = p.plan_id
		INNER JOIN goals g ON g.id = pl.goal_id
		WHERE t.id = ? AND pl.user_id = ?
	`

	var task Task
	err := r.db.QueryRowContext(ctx, query, taskID, userID).Scan(
		&task.GoalID,
		&task.PlanID,
		&task.ID,
		&task.PhaseID,
		&task.GoalTitle,
		&task.PhaseTitle,
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
	)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

// GetSuggestionContextByID 查询任务下一步建议所需上下文。
func (r *Repository) GetSuggestionContextByID(ctx context.Context, userID int64, taskID int64) (appai.TaskSuggestionInput, error) {
	query := `
		SELECT
			g.title,
			pl.title,
			p.title,
			t.title,
			t.description,
			t.status,
			t.priority,
			t.deadline,
			t.phase_id
		FROM tasks t
		INNER JOIN phases p ON p.id = t.phase_id
		INNER JOIN plans pl ON pl.id = p.plan_id
		INNER JOIN goals g ON g.id = pl.goal_id
		WHERE t.id = ? AND pl.user_id = ?
	`

	var input appai.TaskSuggestionInput
	var deadline sql.NullTime
	var phaseID int64
	err := r.db.QueryRowContext(ctx, query, taskID, userID).Scan(
		&input.GoalTitle,
		&input.PlanTitle,
		&input.PhaseTitle,
		&input.Task.Title,
		&input.Task.Description,
		&input.Task.Status,
		&input.Task.Priority,
		&deadline,
		&phaseID,
	)
	if err != nil {
		return appai.TaskSuggestionInput{}, err
	}
	input.Task.Deadline = formatSuggestionDeadline(deadline)

	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT title, description, status, priority, deadline
			FROM tasks
			WHERE phase_id = ? AND id <> ?
			ORDER BY sort_order ASC, id ASC
		`,
		phaseID,
		taskID,
	)
	if err != nil {
		return appai.TaskSuggestionInput{}, err
	}
	defer rows.Close()

	input.SiblingTasks = make([]appai.SuggestionTaskDigest, 0)
	for rows.Next() {
		var task appai.SuggestionTaskDigest
		var siblingDeadline sql.NullTime
		if err := rows.Scan(
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&siblingDeadline,
		); err != nil {
			return appai.TaskSuggestionInput{}, err
		}
		task.Deadline = formatSuggestionDeadline(siblingDeadline)
		input.SiblingTasks = append(input.SiblingTasks, task)
	}
	if err := rows.Err(); err != nil {
		return appai.TaskSuggestionInput{}, err
	}

	return input, nil
}

// SaveSuggestionByID 保存一条任务执行建议。
func (r *Repository) SaveSuggestionByID(ctx context.Context, userID int64, taskID int64, item appai.NextStepSuggestion) error {
	if err := r.ensureOwnedTaskID(ctx, userID, taskID); err != nil {
		return err
	}

	return storedsuggestion.Upsert(ctx, r.db, userID, storedsuggestion.TargetTypeTask, taskID, item)
}

// GetSavedSuggestionByID 查询一条已保存的任务执行建议。
func (r *Repository) GetSavedSuggestionByID(ctx context.Context, userID int64, taskID int64) (appai.NextStepSuggestion, error) {
	if err := r.ensureOwnedTaskID(ctx, userID, taskID); err != nil {
		return appai.NextStepSuggestion{}, err
	}

	return storedsuggestion.Get(ctx, r.db, userID, storedsuggestion.TargetTypeTask, taskID)
}

// ListByUserID 查询当前用户的任务列表。
func (r *Repository) ListByUserID(ctx context.Context, userID int64, req ListTasksRequest) ([]Task, int, error) {
	baseQuery := `
		SELECT
			pl.goal_id,
			p.plan_id,
			t.id,
			t.phase_id,
			g.title,
			p.title,
			t.title,
			t.description,
			t.estimated_days,
			t.deliverables,
			t.deadline,
			t.priority,
			t.status,
			t.sort_order,
			t.created_at,
			t.updated_at
		FROM tasks t
		INNER JOIN phases p ON p.id = t.phase_id
		INNER JOIN plans pl ON pl.id = p.plan_id
		INNER JOIN goals g ON g.id = pl.goal_id
		WHERE pl.user_id = ?
	`
	countQuery := `
		SELECT COUNT(*)
		FROM tasks t
		INNER JOIN phases p ON p.id = t.phase_id
		INNER JOIN plans pl ON pl.id = p.plan_id
		INNER JOIN goals g ON g.id = pl.goal_id
		WHERE pl.user_id = ?
	`

	args := make([]any, 0, 6)
	args = append(args, userID)

	conditions := make([]string, 0, 3)
	if req.Status != "" {
		conditions = append(conditions, "t.status = ?")
		args = append(args, req.Status)
	}
	if req.GoalID > 0 {
		conditions = append(conditions, "pl.goal_id = ?")
		args = append(args, req.GoalID)
	}
	if req.PhaseID > 0 {
		conditions = append(conditions, "t.phase_id = ?")
		args = append(args, req.PhaseID)
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
		countQuery += " AND " + strings.Join(conditions, " AND ")
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	baseQuery += " ORDER BY g.id DESC, p.sort_order ASC, t.sort_order ASC, t.id ASC"
	baseQuery += " LIMIT ? OFFSET ?"
	offset := (req.Page - 1) * req.PageSize
	queryArgs := append(append([]any{}, args...), req.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, baseQuery, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.GoalID,
			&task.PlanID,
			&task.ID,
			&task.PhaseID,
			&task.GoalTitle,
			&task.PhaseTitle,
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
			return nil, 0, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// Create 为当前用户的某个阶段创建任务。
func (r *Repository) Create(ctx context.Context, userID int64, req CreateTaskRequest) (Task, error) {
	phaseID, err := r.getOwnedPhaseID(ctx, userID, req.PhaseID)
	if err != nil {
		return Task{}, err
	}

	sortOrder := req.SortOrder
	if sortOrder <= 0 {
		sortOrder, err = r.nextSortOrderByPhaseID(ctx, phaseID)
		if err != nil {
			return Task{}, err
		}
	}

	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO tasks (phase_id, title, description, estimated_days, deliverables, deadline, priority, status, sort_order) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		phaseID,
		req.Title,
		req.Description,
		req.EstimatedDays,
		req.Deliverables,
		req.Deadline,
		req.Priority,
		TaskStatusTodo,
		sortOrder,
	)
	if err != nil {
		return Task{}, err
	}

	taskID, err := result.LastInsertId()
	if err != nil {
		return Task{}, err
	}

	return r.GetByID(ctx, userID, taskID)
}

// Update 更新当前用户的一条任务。
func (r *Repository) Update(ctx context.Context, userID int64, taskID int64, req UpdateTaskRequest) (Task, error) {
	phaseID, err := r.getOwnedPhaseID(ctx, userID, req.PhaseID)
	if err != nil {
		return Task{}, err
	}

	result, err := r.db.ExecContext(
		ctx,
		`
			UPDATE tasks t
			INNER JOIN phases p ON p.id = t.phase_id
			INNER JOIN plans pl ON pl.id = p.plan_id
			SET
				t.phase_id = ?,
				t.title = ?,
				t.description = ?,
				t.estimated_days = ?,
				t.deliverables = ?,
				t.deadline = ?,
				t.priority = ?,
				t.sort_order = ?
			WHERE t.id = ? AND pl.user_id = ?
		`,
		phaseID,
		req.Title,
		req.Description,
		req.EstimatedDays,
		req.Deliverables,
		req.Deadline,
		req.Priority,
		req.SortOrder,
		taskID,
		userID,
	)
	if err != nil {
		return Task{}, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return Task{}, err
	}
	if affected == 0 {
		return Task{}, sql.ErrNoRows
	}

	return r.GetByID(ctx, userID, taskID)
}

// UpdateStatus 更新当前用户的一条任务状态。
func (r *Repository) UpdateStatus(ctx context.Context, userID int64, taskID int64, status TaskStatus) (Task, error) {
	query := `
		UPDATE tasks t
		INNER JOIN phases p ON p.id = t.phase_id
		INNER JOIN plans pl ON pl.id = p.plan_id
		SET t.status = ?
		WHERE t.id = ? AND pl.user_id = ?
	`

	result, err := r.db.ExecContext(ctx, query, status, taskID, userID)
	if err != nil {
		return Task{}, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return Task{}, err
	}
	if affected == 0 {
		return Task{}, sql.ErrNoRows
	}

	return r.GetByID(ctx, userID, taskID)
}

// Delete 删除当前用户的一条任务。
func (r *Repository) Delete(ctx context.Context, userID int64, taskID int64) error {
	if err := r.ensureOwnedTaskID(ctx, userID, taskID); err != nil {
		return err
	}

	if _, err := r.db.ExecContext(
		ctx,
		`
			DELETE FROM ai_suggestions
			WHERE user_id = ? AND target_type = ? AND target_id = ?
		`,
		userID,
		string(storedsuggestion.TargetTypeTask),
		taskID,
	); err != nil {
		return err
	}

	result, err := r.db.ExecContext(
		ctx,
		`
			DELETE t
			FROM tasks t
			INNER JOIN phases p ON p.id = t.phase_id
			INNER JOIN plans pl ON pl.id = p.plan_id
			WHERE t.id = ? AND pl.user_id = ?
		`,
		taskID,
		userID,
	)
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

// sort 任务排序
func (r *Repository) SortByPhaseID(ctx context.Context, userID int64, phaseID int64, taskIDs []int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 确认这个阶段属于当前用户
	var ownedPhasesID int64
	err = tx.QueryRowContext(
		ctx,
		`SELECT p.id FROM phases p INNER JOIN plans pl ON pl.id = p.plan_id WHERE p.id = ? AND pl.user_id =?`,
		phaseID,
		userID,
	).Scan(&ownedPhasesID)
	if err != nil {
		return err
	}

	rows, err := tx.QueryContext(
		ctx,
		`SELECT id FROM tasks WHERE phase_id = ?`,
		ownedPhasesID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	existing := make(map[int64]struct{})
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return err
		}
		existing[id] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	// 校验每个传入的taskID是否都存在于该阶段
	for _, taskID := range taskIDs {
		if _, ok := existing[taskID]; !ok {
			return sql.ErrNoRows
		}
	}

	// 按顺序更新 sort_order
	for index, taskID := range taskIDs {
		_, err := tx.ExecContext(
			ctx,
			`UPDATE tasks SET sort_order = ? WHERE id = ? AND phase_id = ?`,
			index+1,
			taskID,
			ownedPhasesID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) getOwnedPhaseID(ctx context.Context, userID int64, phaseID int64) (int64, error) {
	var ownedPhaseID int64
	err := r.db.QueryRowContext(
		ctx,
		`
			SELECT p.id
			FROM phases p
			INNER JOIN plans pl ON pl.id = p.plan_id
			WHERE p.id = ? AND pl.user_id = ?
		`,
		phaseID,
		userID,
	).Scan(&ownedPhaseID)
	if err != nil {
		return 0, err
	}

	return ownedPhaseID, nil
}

func (r *Repository) nextSortOrderByPhaseID(ctx context.Context, phaseID int64) (int, error) {
	var sortOrder int
	err := r.db.QueryRowContext(
		ctx,
		`SELECT COALESCE(MAX(sort_order), 0) + 1 FROM tasks WHERE phase_id = ?`,
		phaseID,
	).Scan(&sortOrder)
	if err != nil {
		return 0, err
	}

	return sortOrder, nil
}

func formatSuggestionDeadline(deadline sql.NullTime) string {
	if !deadline.Valid {
		return ""
	}
	return deadline.Time.Format(time.RFC3339)
}

func (r *Repository) ensureOwnedTaskID(ctx context.Context, userID int64, taskID int64) error {
	var ownedTaskID int64
	err := r.db.QueryRowContext(
		ctx,
		`
			SELECT t.id
			FROM tasks t
			INNER JOIN phases p ON p.id = t.phase_id
			INNER JOIN plans pl ON pl.id = p.plan_id
			WHERE t.id = ? AND pl.user_id = ?
		`,
		taskID,
		userID,
	).Scan(&ownedTaskID)
	return err
}
