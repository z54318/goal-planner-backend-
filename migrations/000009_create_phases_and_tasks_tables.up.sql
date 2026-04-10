CREATE TABLE phases (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '阶段ID',
    plan_id BIGINT UNSIGNED NOT NULL COMMENT '所属计划ID',
    title VARCHAR(255) NOT NULL COMMENT '阶段标题',
    description TEXT COMMENT '阶段描述',
    sort_order INT NOT NULL DEFAULT 0 COMMENT '阶段顺序',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    KEY idx_phases_plan_id (plan_id),
    CONSTRAINT fk_phases_plan_id FOREIGN KEY (plan_id) REFERENCES plans(id)
) COMMENT='计划阶段表';

CREATE TABLE tasks (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '任务ID',
    phase_id BIGINT UNSIGNED NOT NULL COMMENT '所属阶段ID',
    title VARCHAR(255) NOT NULL COMMENT '任务标题',
    description TEXT COMMENT '任务描述',
    estimated_days INT NOT NULL DEFAULT 0 COMMENT '预估耗时天数',
    deliverables VARCHAR(255) DEFAULT '' COMMENT '交付物',
    status VARCHAR(50) NOT NULL DEFAULT 'todo' COMMENT '任务状态',
    sort_order INT NOT NULL DEFAULT 0 COMMENT '任务顺序',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    KEY idx_tasks_phase_id (phase_id),
    CONSTRAINT fk_tasks_phase_id FOREIGN KEY (phase_id) REFERENCES phases(id)
) COMMENT='阶段任务表';
