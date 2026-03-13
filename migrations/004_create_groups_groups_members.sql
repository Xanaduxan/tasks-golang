CREATE TABLE IF NOT EXISTS groups
(
    id         UUID PRIMARY KEY,
    name       TEXT      NOT NULL,
    created_by UUID      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS group_members
(
    group_id   UUID      NOT NULL REFERENCES groups (id) ON DELETE CASCADE,
    user_id    UUID      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role       TEXT      NOT NULL DEFAULT 'member',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (group_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_group_members_user_id ON group_members (user_id);