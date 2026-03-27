-- Production Service Migrations

CREATE TABLE IF NOT EXISTS production_activities (
    id BIGSERIAL PRIMARY KEY,
    parcel_id BIGINT NOT NULL,
    activity_type VARCHAR(100) NOT NULL,
    description TEXT,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_activities_parcel_id ON production_activities(parcel_id);
CREATE INDEX IF NOT EXISTS idx_activities_activity_type ON production_activities(activity_type);
CREATE INDEX IF NOT EXISTS idx_activities_timestamp ON production_activities(timestamp);
