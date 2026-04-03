package auth

import (
	"context"
	"database/sql"
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
		SELECT id, username, password_hash, status
		FROM users
		WHERE username = ?
	`

	var user User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Status,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// CreateUser 创建新用户。
func (r *Repository) CreateUser(ctx context.Context, username string, email string, passwordHash string) (int64, error) {
	query := `
		INSERT INTO users (username, email, password_hash, status)
		VALUES (?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query, username, email, passwordHash, "active")
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
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
