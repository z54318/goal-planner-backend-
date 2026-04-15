package rbac

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

var (
	// ErrRoleInUse 表示角色仍在使用中。
	ErrRoleInUse = errors.New("role in use")
	// ErrPermissionInUse 表示权限仍在使用中。
	ErrPermissionInUse = errors.New("permission in use")
	// ErrPermissionNotFound 表示权限不存在。
	ErrPermissionNotFound = errors.New("permission not found")
)

// Repository 负责 rbac 模块和数据库打交道。
type Repository struct {
	db *sql.DB
}

// NewRepository 创建 RBAC 仓库对象。
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// ListRoles 查询角色列表及其权限绑定。
func (r *Repository) ListRoles(ctx context.Context) ([]Role, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, code, created_at, updated_at
		FROM roles
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]Role, 0)
	roleMap := make(map[int64]int)
	roleIDs := make([]int64, 0)
	for rows.Next() {
		var role Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Code, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		roleMap[role.ID] = len(roles)
		roleIDs = append(roleIDs, role.ID)
		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(roleIDs) == 0 {
		return roles, nil
	}

	bindRows, err := r.db.QueryContext(ctx, `
		SELECT role_id, permission_id
		FROM role_permissions
		ORDER BY role_id ASC, permission_id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer bindRows.Close()

	for bindRows.Next() {
		var roleID, permissionID int64
		if err := bindRows.Scan(&roleID, &permissionID); err != nil {
			return nil, err
		}
		index, ok := roleMap[roleID]
		if !ok {
			continue
		}
		roles[index].PermissionIDs = append(roles[index].PermissionIDs, permissionID)
	}
	if err := bindRows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

// GetRoleByID 按ID查询单个角色。
func (r *Repository) GetRoleByID(ctx context.Context, id int64) (Role, error) {
	var role Role
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, code, created_at, updated_at
		FROM roles
		WHERE id = ?
	`, id).Scan(
		&role.ID,
		&role.Name,
		&role.Code,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		return Role{}, err
	}
	return role, nil
}

// CreateRole 新增角色。
func (r *Repository) CreateRole(ctx context.Context, req CreateRoleRequest) (Role, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO roles (name, code)
		VALUES (?, ?)
	`, req.Name, req.Code)
	if err != nil {
		return Role{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Role{}, err
	}

	return r.GetRoleByID(ctx, id)
}

// DeleteRole 删除角色
func (r *Repository) DeleteRole(ctx context.Context, roleID int64) error {
	var exists int
	err := r.db.QueryRowContext(ctx, `SELECT 1 FROM roles WHERE id = ?`, roleID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return err
	}

	var userBindCount int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM user_roles WHERE role_id = ?`, roleID).Scan(&userBindCount); err != nil {
		return err
	}

	if userBindCount > 0 {
		return ErrRoleInUse
	}

	var permissionBindCount int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM role_permissions WHERE role_id = ?`, roleID).Scan(&permissionBindCount); err != nil {
		return err
	}
	if permissionBindCount > 0 {
		return ErrRoleInUse
	}

	result, err := r.db.ExecContext(ctx, `DELETE FROM roles WHERE id = ?`, roleID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// RoleExists 判断角色是否存在。
func (r *Repository) RoleExists(ctx context.Context, roleID int64) (bool, error) {
	var exists int
	err := r.db.QueryRowContext(ctx, `SELECT 1 FROM roles WHERE id = ?`, roleID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ListPermissions 查询全部权限。
func (r *Repository) ListPermissions(ctx context.Context) ([]Permission, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, code, created_at, updated_at
		FROM permissions
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := make([]Permission, 0)
	for rows.Next() {
		var permission Permission
		if err := rows.Scan(&permission.ID, &permission.Name, &permission.Code, &permission.CreatedAt, &permission.UpdatedAt); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetPermissionByID 按ID查询单个权限。
func (r *Repository) GetPermissionByID(ctx context.Context, id int64) (Permission, error) {
	var permission Permission
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, code, created_at, updated_at
		FROM permissions
		WHERE id = ?
	`, id).Scan(
		&permission.ID,
		&permission.Name,
		&permission.Code,
		&permission.CreatedAt,
		&permission.UpdatedAt,
	)
	if err != nil {
		return Permission{}, err
	}
	return permission, nil
}

// CreatePermission 新增权限。
func (r *Repository) CreatePermission(ctx context.Context, req CreatePermissionRequest) (Permission, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO permissions (name, code)
		VALUES (?, ?)
	`, req.Name, req.Code)
	if err != nil {
		return Permission{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Permission{}, err
	}

	return r.GetPermissionByID(ctx, id)
}

// UpdatePermission 更新权限，并同步菜单上的 permission_code。
func (r *Repository) UpdatePermission(ctx context.Context, id int64, req UpdatePermissionRequest) (Permission, error) {
	current, err := r.GetPermissionByID(ctx, id)
	if err != nil {
		return Permission{}, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Permission{}, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `
		UPDATE permissions
		SET name = ?, code = ?
		WHERE id = ?
	`, req.Name, req.Code, id); err != nil {
		return Permission{}, err
	}

	if current.Code != req.Code {
		if _, err := tx.ExecContext(ctx, `
			UPDATE menus
			SET permission_code = ?
			WHERE permission_code = ?
		`, req.Code, current.Code); err != nil {
			return Permission{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return Permission{}, err
	}

	return r.GetPermissionByID(ctx, id)
}

// DeletePermission 删除权限。
func (r *Repository) DeletePermission(ctx context.Context, id int64) error {
	permission, err := r.GetPermissionByID(ctx, id)
	if err != nil {
		return err
	}

	var roleRefCount int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM role_permissions WHERE permission_id = ?`, id).Scan(&roleRefCount); err != nil {
		return err
	}
	if roleRefCount > 0 {
		return ErrPermissionInUse
	}

	var menuRefCount int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM menus WHERE permission_code = ?`, permission.Code).Scan(&menuRefCount); err != nil {
		return err
	}
	if menuRefCount > 0 {
		return ErrPermissionInUse
	}

	result, err := r.db.ExecContext(ctx, `DELETE FROM permissions WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// GetRolePermissionIDs 查询角色已绑定的权限ID。
func (r *Repository) GetRolePermissionIDs(ctx context.Context, roleID int64) ([]int64, error) {
	exists, err := r.RoleExists(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, sql.ErrNoRows
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT permission_id
		FROM role_permissions
		WHERE role_id = ?
		ORDER BY permission_id ASC
	`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

// ReplaceRolePermissions 重置角色的权限绑定。
func (r *Repository) ReplaceRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
	exists, err := r.RoleExists(ctx, roleID)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}

	if len(permissionIDs) > 0 {
		query := `SELECT COUNT(1) FROM permissions WHERE id IN (?` + strings.Repeat(",?", len(permissionIDs)-1) + `)`
		args := make([]any, 0, len(permissionIDs))
		for _, id := range permissionIDs {
			args = append(args, id)
		}
		var count int
		if err := r.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
			return err
		}
		if count != len(permissionIDs) {
			return ErrPermissionNotFound
		}
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM role_permissions WHERE role_id = ?`, roleID); err != nil {
		return err
	}

	for _, permissionID := range permissionIDs {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO role_permissions (role_id, permission_id)
			VALUES (?, ?)
		`, roleID, permissionID); err != nil {
			return err
		}
	}

	return tx.Commit()
}
