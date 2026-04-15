package suggestion

import (
	"context"
	"database/sql"
	"encoding/json"

	appai "goal-planner/internal/infra/ai"
)

// TargetType 表示执行建议所属对象类型。
type TargetType string

const (
	// TargetTypePlan 表示计划执行建议。
	TargetTypePlan TargetType = "plan"
	// TargetTypePhase 表示阶段执行建议。
	TargetTypePhase TargetType = "phase"
	// TargetTypeTask 表示任务执行建议。
	TargetTypeTask TargetType = "task"
)

// Upsert 保存或覆盖一条执行建议。
func Upsert(ctx context.Context, db *sql.DB, userID int64, targetType TargetType, targetID int64, item appai.NextStepSuggestion) error {
	checklistJSON, err := json.Marshal(item.Checklist)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(
		ctx,
		`
			INSERT INTO ai_suggestions (
				user_id,
				target_type,
				target_id,
				summary,
				next_action,
				reason,
				checklist_json,
				risk
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				summary = VALUES(summary),
				next_action = VALUES(next_action),
				reason = VALUES(reason),
				checklist_json = VALUES(checklist_json),
				risk = VALUES(risk),
				updated_at = CURRENT_TIMESTAMP
		`,
		userID,
		targetType,
		targetID,
		item.Summary,
		item.NextAction,
		item.Reason,
		string(checklistJSON),
		item.Risk,
	)
	return err
}

// Get 查询一条已保存的执行建议。
func Get(ctx context.Context, db *sql.DB, userID int64, targetType TargetType, targetID int64) (appai.NextStepSuggestion, error) {
	query := `
		SELECT summary, next_action, reason, checklist_json, risk
		FROM ai_suggestions
		WHERE user_id = ? AND target_type = ? AND target_id = ?
	`

	var item appai.NextStepSuggestion
	var checklistJSON string
	err := db.QueryRowContext(ctx, query, userID, targetType, targetID).Scan(
		&item.Summary,
		&item.NextAction,
		&item.Reason,
		&checklistJSON,
		&item.Risk,
	)
	if err != nil {
		return appai.NextStepSuggestion{}, err
	}

	if checklistJSON != "" {
		if err := json.Unmarshal([]byte(checklistJSON), &item.Checklist); err != nil {
			return appai.NextStepSuggestion{}, err
		}
	}
	if item.Checklist == nil {
		item.Checklist = make([]string, 0)
	}

	return item, nil
}
