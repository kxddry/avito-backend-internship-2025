UPDATE users
SET name = $2,
    is_active = $3,
    team = $4
WHERE id = $1;

