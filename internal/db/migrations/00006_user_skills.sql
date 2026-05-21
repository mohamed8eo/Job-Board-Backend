-- +goose Up
CREATE TABLE user_skills (
  id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
  UNIQUE(user_id, skill_id)
);

-- +goose Down
DROP TABLE user_skills;
