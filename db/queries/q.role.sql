-- name: SetupTenantRole :one
INSERT INTO roles (name, of_store)
VALUES ($1, $2)
RETURNING *;
