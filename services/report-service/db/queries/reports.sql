-- Report queries

-- name: CreateReport :one
INSERT INTO reports (producer_id, report_type, url)
VALUES ($1, $2, $3)
RETURNING id, producer_id, report_type, url, generated_at, created_at;

-- name: GetReport :one
SELECT id, producer_id, report_type, url, generated_at, created_at
FROM reports WHERE id = $1;

-- name: ListReports :many
SELECT id, producer_id, report_type, url, generated_at, created_at
FROM reports WHERE producer_id = $1
ORDER BY generated_at DESC
LIMIT $2 OFFSET $3;
