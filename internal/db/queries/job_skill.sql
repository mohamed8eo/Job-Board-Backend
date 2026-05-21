-- name: AddSkillToJob :one
INSERT INTO job_skills (job_app_id, skill_id)
VALUES ($1, $2)
RETURNING *;

-- name: ListSkillsByJob :many
SELECT s.* FROM skills s
JOIN job_skills js ON js.skill_id = s.id
WHERE js.job_app_id = $1;

-- name: RemoveSkillFromJob :exec
DELETE FROM job_skills
WHERE job_app_id = $1 AND skill_id = $2;
