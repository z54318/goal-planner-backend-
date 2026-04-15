CREATE TABLE ai_suggestions (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '建议ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    target_type VARCHAR(20) NOT NULL COMMENT '对象类型: plan/phase/task',
    target_id BIGINT UNSIGNED NOT NULL COMMENT '对象ID',
    summary VARCHAR(255) NOT NULL COMMENT '建议摘要',
    next_action TEXT NOT NULL COMMENT '建议动作',
    reason TEXT NOT NULL COMMENT '建议原因',
    checklist_json TEXT NOT NULL COMMENT '执行清单JSON',
    risk TEXT NOT NULL COMMENT '主要风险',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_ai_suggestions_target (user_id, target_type, target_id),
    KEY idx_ai_suggestions_lookup (target_type, target_id)
) COMMENT='AI执行建议表';

