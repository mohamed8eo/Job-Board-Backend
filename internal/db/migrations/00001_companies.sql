-- +goose Up
CREATE TABLE companies (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name       TEXT NOT NULL,
  email      TEXT NOT NULL UNIQUE,
  password   TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE companies;
