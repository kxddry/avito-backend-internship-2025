UPDATE pull_requests
SET name = $2,
    author_id = $3,
    status = $4,
    reviewers = $5,
    merged_at = $6
WHERE id = $1
RETURNING created_at, merged_at;

