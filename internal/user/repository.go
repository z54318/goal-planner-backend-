package user

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

var (
	// ErrRoleNotFound 表示存在无效角色ID。
	ErrRoleNotFound = errors.New("role not found")
)

// Repository 负责 user 模块和数据库打交道。
type Repository struct {
	db *sql.DB
}

// NewRepository 创建用户仓库对象。
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// List 查询用户列表及其角色。
func (r *Repository) List(ctx context.Context) ([]User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, username, nickname, email, status, created_at, updated_at
		FROM users
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	userIndex := make(map[int64]int)
	userIDs := make([]int64, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Nickname,
			&user.Email,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}

		userIndex[user.ID] = len(users)
		userIDs = append(userIDs, user.ID)
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(userIDs) == 0 {
		return users, nil
	}

	roleRows, err := r.db.QueryContext(ctx, `
		SELECT ur.user_id, r.id, r.name, r.code
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		ORDER BY ur.user_id ASC, r.id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer roleRows.Close()

	for roleRows.Next() {
		var userID int64
		var role Role
		if err := roleRows.Scan(&userID, &role.ID, &role.Name, &role.Code); err != nil {
			return nil, err
		}
		index, ok := userIndex[userID]
		if !ok {
			continue
		}
		users[index].Roles = append(users[index].Roles, role)
		users[index].RoleIDs = append(users[index].RoleIDs, role.ID)
	}
	if err := roleRows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetByID 按ID查询单个用户及其角色。
func (r *Repository) GetByID(ctx context.Context, userID int64) (User, error) {
	var user User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, username, nickname, email, status, created_at, updated_at
		FROM users
		WHERE id = ?
	`, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Nickname,
		&user.Email,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return User{}, err
	}

	roleRows, err := r.db.QueryContext(ctx, `
		SELECT r.id, r.name, r.code
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_id = ?
		ORDER BY r.id ASC
	`, userID)
	if err != nil {
		return User{}, err
	}
	defer roleRows.Close()

	for roleRows.Next() {
		var role Role
		if err := roleRows.Scan(&role.ID, &role.Name, &role.Code); err != nil {
			return User{}, err
		}
		user.Roles = append(user.Roles, role)
		user.RoleIDs = append(user.RoleIDs, role.ID)
	}
	if err := roleRows.Err(); err != nil {
		return User{}, err
	}

	return user, nil
}

// ReplaceRoles 更新用户角色绑定。
func (r *Repository) ReplaceRoles(ctx context.Context, userID int64, roleIDs []int64) (User, error) {
	if _, err := r.GetByID(ctx, userID); err != nil {
		return User{}, err
	}

	if len(roleIDs) > 0 {
		query := `SELECT COUNT(1) FROM roles WHERE id IN (?` + strings.Repeat(",?", len(roleIDs)-1) + `)`
		args := make([]any, 0, len(roleIDs))
		for _, id := range roleIDs {
			args = append(args, id)
		}

		var count int
		if err := r.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
			return User{}, err
		}
		if count != len(roleIDs) {
			return User{}, ErrRoleNotFound
		}
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return User{}, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM user_roles WHERE user_id = ?`, userID); err != nil {
		return User{}, err
	}

	for _, roleID := range roleIDs {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO user_roles (user_id, role_id)
			VALUES (?, ?)
		`, userID, roleID); err != nil {
			return User{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return User{}, err
	}

	return r.GetByID(ctx, userID)
}
