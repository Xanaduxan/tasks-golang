CREATE TABLE IF NOT EXISTS task_user_analytics
(
    user_id         UUID PRIMARY KEY,
    tasks_completed BIGINT      NOT NULL DEFAULT 0,
    last_event_at   TIMESTAMPTZ NULL
);