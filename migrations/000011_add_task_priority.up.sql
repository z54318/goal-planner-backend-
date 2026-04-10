ALTER TABLE tasks
ADD COLUMN priority VARCHAR(20) NOT NULL DEFAULT 'medium' COMMENT '任务优先级：high高 medium中 low低' AFTER deliverables;
