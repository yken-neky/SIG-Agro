-- Parcel queries

-- name: CreateParcel :one
INSERT INTO parcels (producer_id, name, description, geometry, area_hectares, crop_type)
VALUES ($1, $2, $3, ST_GeomFromText($4, 4326), $5, $6)
RETURNING id, producer_id, name, description, ST_AsText(geometry) as geometry, area_hectares, crop_type, created_at;

-- name: GetParcel :one
SELECT id, producer_id, name, description, ST_AsText(geometry) as geometry, area_hectares, crop_type, created_at
FROM parcels WHERE id = $1;

-- name: ListParcels :many
SELECT id, producer_id, name, description, ST_AsText(geometry) as geometry, area_hectares, crop_type, created_at
FROM parcels WHERE producer_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateParcel :exec
UPDATE parcels
SET name = $2, description = $3, crop_type = $4, geometry = ST_GeomFromText($5, 4326), updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteParcel :exec
DELETE FROM parcels WHERE id = $1;

-- name: QueryByGeometry :many
SELECT id, producer_id, name, description, ST_AsText(geometry) as geometry, area_hectares, crop_type, created_at
FROM parcels
WHERE ST_Intersects(geometry, ST_GeomFromText($1, 4326))
ORDER BY created_at DESC;

-- name: GetParcelCount :one
SELECT COUNT(*) FROM parcels WHERE producer_id = $1;
