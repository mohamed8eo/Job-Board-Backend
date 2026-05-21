-- name: CreateSkill :one
INSERT INTO skills (name)
VALUES ($1)
RETURNING *;

-- name: GetSkillByID :one
SELECT * FROM skills
WHERE id = $1;

-- name: GetSkillByName :one
SELECT * FROM skills
WHERE name = $1;

-- name: ListSkills :many
SELECT * FROM skills;

-- name: DeleteSkill :exec
DELETE FROM skills
WHERE id = $1;
