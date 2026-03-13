ALTER TABLE tasks
ADD CONSTRAINT tasks_status_check
CHECK (status IN ('created', 'in_progress', 'done'));