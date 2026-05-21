-- name: AddSkillToUser :one
INSERT INTO user_skills (user_id, skill_id)
VALUES ($1, $2)
RETURNING *;

-- name: ListSkillsByUser :many
SELECT s.* FROM skills s
JOIN user_skills us ON us.skill_id = s.id
WHERE us.user_id = $1;

-- name: RemoveSkillFromUser :exec
DELETE FROM user_skills
WHERE user_id = $1 AND skill_id = $2;
