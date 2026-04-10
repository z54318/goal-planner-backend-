CREATE TABLE plans (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '计划ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    goal_id BIGINT UNSIGNED NOT NULL COMMENT '目标ID',
    title VARCHAR(255) NOT NULL COMMENT '计划标题',
    overview TEXT COMMENT '计划概述',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_plans_goal_id (goal_id),
    KEY idx_plans_user_id (user_id),
    CONSTRAINT fk_plans_user_id FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_plans_goal_id FOREIGN KEY (goal_id) REFERENCES goals(id)
) COMMENT='计划表';
