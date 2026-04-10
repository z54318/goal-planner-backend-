UPDATE menus
SET
    path = '/goals/list',
    component = 'goals/list',
    permission_code = 'goal:list'
WHERE id = 2;

UPDATE menus
SET hidden = 1
WHERE id IN (3, 4, 5);
