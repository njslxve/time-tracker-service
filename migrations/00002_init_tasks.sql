-- +goose Up
CREATE TABLE IF NOT EXISTS tasks(
  id uuid PRIMARY KEY,
  task_id TEXT NOT NULL,
  user_id integer NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  start_time timestamptz NOT NULL,
  end_time timestamptz,
  duration integer
);

-- +goose Down
DROP TABLE tasks;