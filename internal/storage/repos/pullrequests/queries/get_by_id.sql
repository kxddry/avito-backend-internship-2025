SELECT id,
       name,
       author_id,
       status,
       reviewers,
       created_at,
       merged_at
FROM pull_requests
WHERE id = $1;

