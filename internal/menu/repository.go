package menu

import (
	"context"
	"database/sql"
	"errors"
)

var (
	// ErrMenuHasChildren 表示当前菜单下仍有子菜单。
	ErrMenuHasChildren = errors.New("menu has children")
)

// Repository 负责 menu 模块和数据库打交道。
type Repository struct {
	db *sql.DB
}

// NewRepository 创建菜单仓库对象。
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// List 查询全部菜单。
func (r *Repository) List(ctx context.Context) ([]Menu, error) {
	query := `
		SELECT id, parent_id, name, path, component, icon, sort_order, permission_code, hidden, created_at, updated_at
		FROM menus
		ORDER BY parent_id ASC, sort_order ASC, id ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
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
			&menu.CreatedAt,
			&menu.UpdatedAt,
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

// GetByID 按ID查询单个菜单。
func (r *Repository) GetByID(ctx context.Context, id int64) (Menu, error) {
	query := `
		SELECT id, parent_id, name, path, component, icon, sort_order, permission_code, hidden, created_at, updated_at
		FROM menus
		WHERE id = ?
	`

	var menu Menu
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&menu.ID,
		&menu.ParentID,
		&menu.Name,
		&menu.Path,
		&menu.Component,
		&menu.Icon,
		&menu.SortOrder,
		&menu.PermissionCode,
		&menu.Hidden,
		&menu.CreatedAt,
		&menu.UpdatedAt,
	)
	if err != nil {
		return Menu{}, err
	}

	return menu, nil
}

// ParentExists 判断父菜单是否存在。
func (r *Repository) ParentExists(ctx context.Context, parentID int64) (bool, error) {
	if parentID == 0 {
		return true, nil
	}

	var exists int
	err := r.db.QueryRowContext(ctx, `SELECT 1 FROM menus WHERE id = ?`, parentID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// Create 新增一条菜单记录。
func (r *Repository) Create(ctx context.Context, req CreateMenuRequest) (Menu, error) {
	query := `
		INSERT INTO menus (parent_id, name, path, component, icon, sort_order, permission_code, hidden)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		req.ParentID,
		req.Name,
		req.Path,
		req.Component,
		req.Icon,
		req.SortOrder,
		req.PermissionCode,
		req.Hidden,
	)
	if err != nil {
		return Menu{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Menu{}, err
	}

	return r.GetByID(ctx, id)
}

// Update 更新一条菜单记录。
func (r *Repository) Update(ctx context.Context, id int64, req UpdateMenuRequest) (Menu, error) {
	if _, err := r.GetByID(ctx, id); err != nil {
		return Menu{}, err
	}

	query := `
		UPDATE menus
		SET parent_id = ?, name = ?, path = ?, component = ?, icon = ?, sort_order = ?, permission_code = ?, hidden = ?
		WHERE id = ?
	`

	if _, err := r.db.ExecContext(
		ctx,
		query,
		req.ParentID,
		req.Name,
		req.Path,
		req.Component,
		req.Icon,
		req.SortOrder,
		req.PermissionCode,
		req.Hidden,
		id,
	); err != nil {
		return Menu{}, err
	}

	return r.GetByID(ctx, id)
}

// Delete 删除一条菜单记录。
func (r *Repository) Delete(ctx context.Context, id int64) error {
	var childCount int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM menus WHERE parent_id = ?`, id).Scan(&childCount); err != nil {
		return err
	}
	if childCount > 0 {
		return ErrMenuHasChildren
	}

	result, err := r.db.ExecContext(ctx, `DELETE FROM menus WHERE id = ?`, id)
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
