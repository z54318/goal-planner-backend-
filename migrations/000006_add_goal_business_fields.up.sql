ALTER TABLE goals
    ADD COLUMN user_id BIGINT UNSIGNED NULL COMMENT '所属用户ID' AFTER id,
    ADD COLUMN category VARCHAR(100) NOT NULL DEFAULT '' COMMENT '目标分类' AFTER description,
    ADD COLUMN target_deadline DATETIME NULL COMMENT '目标截止时间' AFTER category;

ALTER TABLE goals
    ADD INDEX idx_goals_user_id (user_id);

ALTER TABLE goals
    ADD CONSTRAINT fk_goals_user_id
    FOREIGN KEY (user_id) REFERENCES users(id);
