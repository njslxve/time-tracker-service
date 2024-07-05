-- +goose Up
CREATE INDEX IF NOT EXISTS idx_tasks_task_id_user_id ON tasks(task_id, user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_tasks_task_id_user_id;