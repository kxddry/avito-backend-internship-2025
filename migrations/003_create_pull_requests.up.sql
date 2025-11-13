CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE IF NOT EXISTS pull_requests (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    author_id TEXT REFERENCES users(id),
    status pr_status NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reviewers TEXT[] NOT NULL DEFAULT '{}', -- такое решение я принял из-за ограничения в 2 ревьюера на задачу.
    merged_at TIMESTAMPTZ
    -- Было бы логичнее сделать отдельную таблицу для ревьюеров при большем количестве ревьюеров.
);