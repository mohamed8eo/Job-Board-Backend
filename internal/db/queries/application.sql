-- name: CreateApplication :one
INSERT INTO applications (user_id, job_app_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetApplicationByID :one
SELECT * FROM applications
WHERE id = $1;

-- name: ListApplicationsByJob :many
SELECT * FROM applications
WHERE job_app_id = $1;

-- name: ListApplicationsByUser :many
SELECT * FROM applications
WHERE user_id = $1;

-- name: UpdateApplicationStatus :one
UPDATE applications
SET status = $2
WHERE id = $1
RETURNING *;

-- name: DeleteApplication :exec
DELETE FROM applications
WHERE id = $1;
