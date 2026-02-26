-- name: SetupTenantUser :one
INSERT INTO users (username, email, password, name, of_store, of_role)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByAny :one
SELECT * FROM users WHERE email = $1
    OR username = $2
    OR phone_number = $3
    LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: RegisterUser :one
INSERT INTO users (email, password, name, username, phone_number)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
