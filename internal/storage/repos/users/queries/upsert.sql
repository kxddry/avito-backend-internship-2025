INSERT INTO users (id, name, is_active, team)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE
SET name = EXCLUDED.name,
    is_active = EXCLUDED.is_active,
    team = EXCLUDED.team;

