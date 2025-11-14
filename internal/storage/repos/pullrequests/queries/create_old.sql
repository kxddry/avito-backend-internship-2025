INSERT INTO pull_requests (id, name, author_id, status, reviewers, merged_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING created_at, merged_at;

