ALTER TABLE tasks
    ADD COLUMN IF NOT EXISTS search_vector tsvector;

UPDATE tasks
SET search_vector = to_tsvector('simple', coalesce(name, ''));

CREATE INDEX IF NOT EXISTS idx_tasks_search_vector
    ON tasks USING GIN (search_vector);

CREATE OR REPLACE FUNCTION tasks_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
            to_tsvector('simple', coalesce(NEW.name, ''));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_tasks_search_vector_update ON tasks;

CREATE TRIGGER trg_tasks_search_vector_update
    BEFORE INSERT OR UPDATE OF name
    ON tasks
    FOR EACH ROW
EXECUTE FUNCTION tasks_search_vector_update();