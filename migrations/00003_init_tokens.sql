-- +goose Up
CREATE TABLE IF NOT EXISTS pagination_tokens(
  token TEXT PRIMARY KEY,
  old_limit integer NOT NULL,
  filter_params JSONB,
  is_alive boolean DEFAULT true,
  created_at timestamptz,
  ttl INTERVAL DEFAULT '1 day'
);

-- +goose Down
DROP TABLE pagination_tokens;