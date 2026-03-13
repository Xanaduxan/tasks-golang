ALTER TABLE tasks
ADD COLUMN group_id UUID REFERENCES groups(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_tasks_group_id ON tasks(group_id);