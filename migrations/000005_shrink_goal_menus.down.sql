UPDATE menus
SET
    path = '/goals',
    component = '',
    permission_code = ''
WHERE id = 2;

UPDATE menus
SET hidden = CASE
    WHEN id = 3 THEN 0
    WHEN id = 4 THEN 0
    WHEN id = 5 THEN 1
END
WHERE id IN (3, 4, 5);
