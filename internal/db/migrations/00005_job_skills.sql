-- +goose Up
CREATE TABLE job_skills (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  job_app_id UUID NOT NULL REFERENCES job_apps(id) ON DELETE CASCADE,
  skill_id   UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
  UNIQUE(job_app_id, skill_id)
);

-- +goose Down
DROP TABLE job_skills;
