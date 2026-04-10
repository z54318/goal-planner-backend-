ALTER TABLE users
    ADD COLUMN nickname VARCHAR(100) NOT NULL DEFAULT '' COMMENT '昵称' AFTER username;

UPDATE users
SET nickname = username
WHERE nickname = '';
