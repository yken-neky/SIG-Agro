-- Parcel Service Migrations

-- Enable PostGIS extension
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS postgis_topology;

CREATE TABLE IF NOT EXISTS parcels (
    id BIGSERIAL PRIMARY KEY,
    producer_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    geometry GEOMETRY(POLYGON, 4326) NOT NULL,
    area_hectares NUMERIC(10, 2),
    crop_type VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_parcels_producer_id ON parcels(producer_id);
CREATE INDEX IF NOT EXISTS idx_parcels_geometry ON parcels USING GIST(geometry);
CREATE INDEX IF NOT EXISTS idx_parcels_crop_type ON parcels(crop_type);
