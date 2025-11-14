SELECT
    COUNT(*) as total,
    COUNT(CASE WHEN status = 'OPEN' THEN 1 END) as open,
    COUNT(CASE WHEN status = 'MERGED' THEN 1 END) as merged,
    COUNT(CASE WHEN array_length(reviewers, 1) = 0 OR reviewers IS NULL THEN 1 END) as with0,
    COUNT(CASE WHEN array_length(reviewers, 1) = 1 THEN 1 END) as with1,
    COUNT(CASE WHEN array_length(reviewers, 1) = 2 THEN 1 END) as with2
FROM pull_requests;