CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE IF NOT EXISTS pull_requests (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    author_id UUID REFERENCES users(id),
    status pr_status NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reviewers TEXT[] NOT NULL DEFAULT '{}'
);