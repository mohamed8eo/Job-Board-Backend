-- name: CreateJobApp :one
INSERT INTO job_apps (company_id, title, description, location, remote, salary)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetJobAppByID :one
SELECT * FROM job_apps
WHERE id = $1;

-- name: ListJobApps :many
SELECT * FROM job_apps;

-- name: ListJobAppsByCompany :many
SELECT * FROM job_apps
WHERE company_id = $1;

-- name: UpdateJobApp :one
UPDATE job_apps
SET title = $2, description = $3, location = $4, remote = $5, salary = $6
WHERE id = $1
RETURNING *;

-- name: DeleteJobApp :exec
DELETE FROM job_apps
WHERE id = $1;
