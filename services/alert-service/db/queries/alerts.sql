-- Alert queries

-- name: CreateAlert :one
INSERT INTO alerts (parcel_id, alert_type, severity, message)
VALUES ($1, $2, $3, $4)
RETURNING id, parcel_id, alert_type, severity, message, resolved, created_at;

-- name: GetAlert :one
SELECT id, parcel_id, alert_type, severity, message, resolved, created_at
FROM alerts WHERE id = $1;

-- name: ListAlerts :many
SELECT id, parcel_id, alert_type, severity, message, resolved, created_at
FROM alerts WHERE parcel_id = $1 AND severity = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: ResolveAlert :exec
UPDATE alerts
SET resolved = TRUE, resolved_at = CURRENT_TIMESTAMP
WHERE id = $1;
