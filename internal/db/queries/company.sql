-- name: CreateCompany :one
INSERT INTO companies (name, email, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCompanyByID :one
SELECT * FROM companies
WHERE id = $1;

-- name: GetCompanyByEmail :one
SELECT * FROM companies
WHERE email = $1;

-- name: ListCompanies :many
SELECT * FROM companies;

-- name: DeleteCompany :exec
DELETE FROM companies
WHERE id = $1;
