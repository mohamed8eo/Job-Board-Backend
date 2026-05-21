-- +goose Up
CREATE TABLE applications (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  job_app_id UUID NOT NULL REFERENCES job_apps(id) ON DELETE CASCADE,
  status     TEXT NOT NULL DEFAULT 'pending',
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  UNIQUE(user_id, job_app_id)
);

-- +goose Down
DROP TABLE applications;
