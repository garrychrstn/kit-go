-- name: GetStore :one
SELECT * FROM stores WHERE id = $1;

-- name: GetStoreByName :one
SELECT * FROM stores WHERE name = $1;

-- name: ListStores :many
SELECT * FROM stores ORDER BY created_at DESC;

-- name: CreateStore :one
INSERT INTO stores(name, description, logo, address, phone, category, contacts, of_owner, term_and_service, coordinate, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;
