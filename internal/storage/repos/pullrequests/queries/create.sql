WITH author AS (
    SELECT id AS author_id, team
    FROM users
    WHERE id = $1
),
     team_members AS (
         SELECT u.id
         FROM users u
                  JOIN author a ON u.team = a.team
         WHERE u.id <> a.author_id
           AND u.is_active = TRUE
     ),
     picked AS (
         SELECT COALESCE(
                    -- максимум 2 ревьюера, случайно
                        (array_agg(id ORDER BY random()))[1:2],
                        '{}'::text[]
                ) AS reviewers
         FROM team_members
     ),
     inserted AS (
         INSERT INTO pull_requests (id, name, author_id, reviewers)
             SELECT
                 $2 AS id,
                 $3 AS name,
                 a.author_id,
                 p.reviewers
             FROM author a
                      CROSS JOIN picked p
             RETURNING id, name, author_id, status, created_at, reviewers, merged_at
     )
SELECT *
FROM inserted;
