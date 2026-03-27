-- User queries

-- name: CreateUser :one
INSERT INTO users (email, password_hash, full_name, phone)
VALUES ($1, $2, $3, $4)
RETURNING id, email, full_name, phone, created_at;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, full_name, phone, created_at FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, email, full_name, phone, created_at FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, email, full_name, phone, created_at FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUser :exec
UPDATE users
SET full_name = $2, phone = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: CreateUserRole :exec
INSERT INTO user_roles (user_id, role) VALUES ($1, $2)
ON CONFLICT (user_id, role) DO NOTHING;

-- name: GetUserRoles :many
SELECT role FROM user_roles WHERE user_id = $1;

-- name: UpsertToken :exec
INSERT INTO tokens (user_id, token_hash, expires_at)
VALUES ($1, $2, $3)
ON CONFLICT (token_hash) DO NOTHING;

-- name: GetUserByToken :one
SELECT id, email, full_name FROM users
WHERE id IN (
  SELECT user_id FROM tokens
  WHERE token_hash = $1 AND expires_at > CURRENT_TIMESTAMP
);
