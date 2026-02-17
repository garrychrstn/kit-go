-- name: SetupTenantUser :one
INSERT INTO users (username, email, password, name, of_store, of_role)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;
