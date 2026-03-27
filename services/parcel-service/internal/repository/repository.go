package repository

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

type Database struct {
	*sql.DB
}

func NewDatabase(dsn string) (*Database, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{db}, nil
}

func NewRepository(db *Database) *Repository {
	return &Repository{db: db.DB}
}

func (r *Repository) CreateParcel(ctx context.Context, producerID int64, name, description, geometryWKT string, areaHectares float64, cropType string) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO parcels (producer_id, name, description, geometry, area_hectares, crop_type) VALUES ($1, $2, $3, ST_GeomFromText($4, 4326), $5, $6) RETURNING id",
		producerID, name, description, geometryWKT, areaHectares, cropType,
	).Scan(&id)
	return id, err
}

func (r *Repository) GetParcel(ctx context.Context, id int64) (int64, string, string, string, float64, string, error) {
	var producerID int64
	var name, description, geometryWKT string
	var areaHectares float64
	var cropType string
	err := r.db.QueryRowContext(ctx,
		"SELECT producer_id, name, description, ST_AsText(geometry), area_hectares, crop_type FROM parcels WHERE id = $1",
		id,
	).Scan(&producerID, &name, &description, &geometryWKT, &areaHectares, &cropType)
	return producerID, name, description, geometryWKT, areaHectares, cropType, err
}

func (r *Repository) ListParcels(ctx context.Context, producerID int64, limit, offset int32) ([]map[string]interface{}, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, producer_id, name, description, ST_AsText(geometry), area_hectares, crop_type, created_at FROM parcels WHERE producer_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3",
		producerID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcels []map[string]interface{}
	for rows.Next() {
		var id, pid int64
		var name, description, geometryWKT, cropType string
		var areaHectares float64
		var createdAt int64
		if err := rows.Scan(&id, &pid, &name, &description, &geometryWKT, &areaHectares, &cropType, &createdAt); err != nil {
			return nil, err
		}
		parcels = append(parcels, map[string]interface{}{
			"id":            id,
			"producer_id":   pid,
			"name":          name,
			"description":   description,
			"geometry_wkt":  geometryWKT,
			"area_hectares": areaHectares,
			"crop_type":     cropType,
			"created_at":    createdAt,
		})
	}
	return parcels, rows.Err()
}

func (r *Repository) UpdateParcel(ctx context.Context, id int64, name, description, cropType, geometryWKT string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE parcels SET name = $2, description = $3, crop_type = $4, geometry = ST_GeomFromText($5, 4326), updated_at = CURRENT_TIMESTAMP WHERE id = $1",
		id, name, description, cropType, geometryWKT,
	)
	return err
}

func (r *Repository) DeleteParcel(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM parcels WHERE id = $1", id)
	return err
}

func (r *Repository) QueryByGeometry(ctx context.Context, geometryWKT string) ([]map[string]interface{}, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, producer_id, name, description, ST_AsText(geometry), area_hectares, crop_type, created_at FROM parcels WHERE ST_Intersects(geometry, ST_GeomFromText($1, 4326)) ORDER BY created_at DESC",
		geometryWKT,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcels []map[string]interface{}
	for rows.Next() {
		var id, pid int64
		var name, description, geometryWKT, cropType string
		var areaHectares float64
		var createdAt int64
		if err := rows.Scan(&id, &pid, &name, &description, &geometryWKT, &areaHectares, &cropType, &createdAt); err != nil {
			return nil, err
		}
		parcels = append(parcels, map[string]interface{}{
			"id":            id,
			"producer_id":   pid,
			"name":          name,
			"description":   description,
			"geometry_wkt":  geometryWKT,
			"area_hectares": areaHectares,
			"crop_type":     cropType,
			"created_at":    createdAt,
		})
	}
	return parcels, rows.Err()
}

func (r *Repository) GetParcelCount(ctx context.Context, producerID int64) (int32, error) {
	var count int32
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM parcels WHERE producer_id = $1", producerID).Scan(&count)
	return count, err
}
