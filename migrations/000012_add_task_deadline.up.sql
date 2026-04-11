ALTER TABLE tasks
    ADD COLUMN deadline DATETIME NULL COMMENT '任务截止时间' AFTER deliverables;
