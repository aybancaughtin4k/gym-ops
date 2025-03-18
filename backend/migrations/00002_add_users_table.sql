-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  fullname TEXT NOT NULL,
  email citext UNIQUE NOT NULL,
  password BYTEA NOT NULL,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
