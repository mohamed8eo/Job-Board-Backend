-- +goose Up
CREATE TABLE skills (
  id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE skills;
