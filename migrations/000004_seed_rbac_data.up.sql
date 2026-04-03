INSERT INTO roles (name, code) VALUES
('管理员', 'admin'),
('普通用户', 'user');

INSERT INTO permissions (name, code) VALUES
('查看仪表盘', 'dashboard:view'),
('查看目标列表', 'goal:list'),
('创建目标', 'goal:create'),
('查看目标详情', 'goal:detail');

INSERT INTO menus (id, parent_id, name, path, component, icon, sort_order, permission_code, hidden) VALUES
(1, 0, '仪表盘', '/dashboard', 'dashboard/index', 'DashboardOutlined', 1, 'dashboard:view', 0),
(2, 0, '目标管理', '/goals', '', 'FlagOutlined', 2, '', 0),
(3, 2, '目标列表', '/goals/list', 'goals/list', '', 1, 'goal:list', 0),
(4, 2, '新建目标', '/goals/create', 'goals/create', '', 2, 'goal:create', 0),
(5, 2, '目标详情', '/goals/detail', 'goals/detail', '', 3, 'goal:detail', 1);

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p
WHERE r.code = 'admin';

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p
WHERE r.code = 'user'
  AND p.code IN ('dashboard:view', 'goal:list', 'goal:create', 'goal:detail');
