-- +goose Up
CREATE TABLE job_apps (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  company_id  UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
  title       TEXT NOT NULL,
  description TEXT NOT NULL,
  location    TEXT NOT NULL,
  remote      BOOLEAN NOT NULL DEFAULT false,
  salary      INT,
  created_at  TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE job_apps;
