SELECT id, name, is_active
FROM users
WHERE team = $1
ORDER BY name ASC;

