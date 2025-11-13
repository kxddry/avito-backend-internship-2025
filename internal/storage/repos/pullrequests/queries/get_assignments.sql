SELECT id,
       name,
       author_id,
       status
FROM pull_requests
WHERE $1 = ANY(reviewers)
ORDER BY created_at DESC;

