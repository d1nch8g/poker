-- name: CreateUser :one
INSERT INTO users (
        name,
        email,
        verified,
        is_admin,
        encrypted_wallet,
        passwhash
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
-- name: ListUsers :many
SELECT *
FROM users;
-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1;