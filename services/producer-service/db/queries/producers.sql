-- Producer queries

-- name: CreateProducer :one
INSERT INTO producers (user_id, name, document_id, phone, email, address)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, name, document_id, phone, email, address, created_at;

-- name: GetProducer :one
SELECT id, user_id, name, document_id, phone, email, address, created_at
FROM producers WHERE id = $1;

-- name: ListProducers :many
SELECT id, user_id, name, document_id, phone, email, address, created_at
FROM producers WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateProducer :exec
UPDATE producers
SET name = $2, phone = $3, email = $4, address = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteProducer :exec
DELETE FROM producers WHERE id = $1;

-- name: GetProducerCount :one
SELECT COUNT(*) FROM producers WHERE user_id = $1;
