ALTER TABLE goals
    DROP FOREIGN KEY fk_goals_user_id;

ALTER TABLE goals
    DROP INDEX idx_goals_user_id,
    DROP COLUMN target_deadline,
    DROP COLUMN category,
    DROP COLUMN user_id;
