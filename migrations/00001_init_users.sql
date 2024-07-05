-- +goose Up
CREATE TABLE IF NOT EXISTS users(
  id uuid PRIMARY KEY,
  user_id SERIAL UNIQUE NOT NULL,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  patronymic TEXT,
  passport TEXT NOT NULL UNIQUE,
  adress TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;