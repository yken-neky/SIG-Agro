-- Production queries

-- name: RecordActivity :one
INSERT INTO production_activities (parcel_id, activity_type, description, timestamp, metadata)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, parcel_id, activity_type, description, timestamp;

-- name: GetActivity :one
SELECT id, parcel_id, activity_type, description, timestamp
FROM production_activities WHERE id = $1;

-- name: ListActivities :many
SELECT id, parcel_id, activity_type, description, timestamp
FROM production_activities WHERE parcel_id = $1
ORDER BY timestamp DESC
LIMIT $2 OFFSET $3;
