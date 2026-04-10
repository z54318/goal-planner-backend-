package auth

import (
	"context"
	"database/sql"
	"errors"
)

// Repository 负责 auth 模块和数据库交互。
type Repository struct {
	db *sql.DB
}

// NewRepository 创建认证仓库对象。
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// GetByUsername 按用户名查询用户。
func (r *Repository) GetByUsername(ctx context.Context, username string) (User, error) {
	query := `
		SELECT id, username, nickname, password_hash, status
		FROM users
		WHERE username = ?
	`

	var user User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Nickname,
		&user.PasswordHash,
		&user.Status,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// CreateUser 创建新用户，并为其绑定默认 user 角色。
func (r *Repository) CreateUser(ctx context.Context, username string, nickname string, email string, passwordHash string) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := `
		INSERT INTO users (username, nickname, email, password_hash, status)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := tx.ExecContext(ctx, query, username, nickname, email, passwordHash, "active")
	if err != nil {
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	var roleID int64
	err = tx.QueryRowContext(ctx, `SELECT id FROM roles WHERE code = ?`, "user").Scan(&roleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("默认角色 user 不存在")
		}
		return 0, err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO user_roles (user_id, role_id)
		VALUES (?, ?)
	`, userID, roleID)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return userID, nil
}

// ListMenusByUserID 查询当前用户可见的菜单列表。
func (r *Repository) ListMenusByUserID(ctx context.Context, userID int64) ([]Menu, error) {
	query := `
		SELECT DISTINCT
			m.id,
			m.parent_id,
			m.name,
			m.path,
			m.component,
			m.icon,
			m.sort_order,
			m.permission_code,
			m.hidden
		FROM menus m
		WHERE m.hidden = 0
		  AND (
			m.permission_code = ''
			OR m.permission_code IN (
				SELECT p.code
				FROM user_roles ur
				JOIN role_permissions rp ON ur.role_id = rp.role_id
				JOIN permissions p ON rp.permission_id = p.id
				WHERE ur.user_id = ?
			)
		  )
		ORDER BY m.parent_id ASC, m.sort_order ASC, m.id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	menus := make([]Menu, 0)
	for rows.Next() {
		var menu Menu
		if err := rows.Scan(
			&menu.ID,
			&menu.ParentID,
			&menu.Name,
			&menu.Path,
			&menu.Component,
			&menu.Icon,
			&menu.SortOrder,
			&menu.PermissionCode,
			&menu.Hidden,
		); err != nil {
			return nil, err
		}

		menus = append(menus, menu)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return menus, nil
}
