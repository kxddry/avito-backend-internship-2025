SELECT
    u.id,
    u.name,
    COALESCE(u.team, ''),
    u.is_active,
    COUNT(DISTINCT CASE WHEN pr.id IS NOT NULL THEN pr.id END) as total_reviews,
    COUNT(DISTINCT CASE WHEN pr.status = 'OPEN' THEN pr.id END) as open_reviews,
    COUNT(DISTINCT CASE WHEN pr.status = 'MERGED' THEN pr.id END) as merged_reviews
FROM users u
         LEFT JOIN pull_requests pr ON u.id = ANY(pr.reviewers)
GROUP BY u.id, u.name, u.team, u.is_active
ORDER BY u.name;