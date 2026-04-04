CREATE TABLE IF NOT EXISTS task_event_log
(
    event_id        UUID PRIMARY KEY,
    task_id         UUID        NOT NULL,
    user_id         UUID        NOT NULL,
    group_id        UUID        NULL,
    event_type      TEXT        NOT NULL,
    status          TEXT        NOT NULL,
    previous_status TEXT        NULL,
    created_at      TIMESTAMPTZ NOT NULL
);