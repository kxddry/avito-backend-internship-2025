SELECT
    t.name,
    COUNT(DISTINCT u.id) as members_total,
    COUNT(DISTINCT CASE WHEN u.is_active THEN u.id END) as members_active,
    COUNT(DISTINCT pr.id) as prs_created,
    COUNT(DISTINCT CASE WHEN pr.status = 'OPEN' THEN pr.id END) as prs_open
FROM teams t
         LEFT JOIN users u ON t.name = u.team
         LEFT JOIN pull_requests pr ON u.id = pr.author_id
GROUP BY t.name
ORDER BY t.name
