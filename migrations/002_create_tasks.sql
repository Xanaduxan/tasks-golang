CREATE TABLE IF NOT EXISTS tasks
(
    id         UUID PRIMARY KEY,
    name       TEXT      NOT NULL,
    deadline   TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    user_id    UUID      NOT NULL REFERENCES users (id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks (user_id);