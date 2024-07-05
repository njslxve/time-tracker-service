-- +goose Up
CREATE INDEX IF NOT EXISTS idx_users_user_id ON users(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_users_user_id;