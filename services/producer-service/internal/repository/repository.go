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

func (r *Repository) CreateProducer(ctx context.Context, userID int64, name, documentID, phone, email, address string) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO producers (user_id, name, document_id, phone, email, address) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		userID, name, documentID, phone, email, address,
	).Scan(&id)
	return id, err
}

func (r *Repository) GetProducer(ctx context.Context, id int64) (int64, string, string, string, string, string, error) {
	var userID int64
	var name, documentID, phone, email, address string
	err := r.db.QueryRowContext(ctx,
		"SELECT user_id, name, document_id, phone, email, address FROM producers WHERE id = $1",
		id,
	).Scan(&userID, &name, &documentID, &phone, &email, &address)
	return userID, name, documentID, phone, email, address, err
}

func (r *Repository) ListProducers(ctx context.Context, userID int64, limit, offset int32) ([]map[string]interface{}, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, user_id, name, document_id, phone, email, address, created_at FROM producers WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3",
		userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var producers []map[string]interface{}
	for rows.Next() {
		var id, uid int64
		var name, documentID, phone, email, address string
		var createdAt int64
		if err := rows.Scan(&id, &uid, &name, &documentID, &phone, &email, &address, &createdAt); err != nil {
			return nil, err
		}
		producers = append(producers, map[string]interface{}{
			"id":          id,
			"user_id":     uid,
			"name":        name,
			"document_id": documentID,
			"phone":       phone,
			"email":       email,
			"address":     address,
			"created_at":  createdAt,
		})
	}
	return producers, rows.Err()
}

func (r *Repository) UpdateProducer(ctx context.Context, id int64, name, phone, email, address string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE producers SET name = $2, phone = $3, email = $4, address = $5, updated_at = CURRENT_TIMESTAMP WHERE id = $1",
		id, name, phone, email, address,
	)
	return err
}

func (r *Repository) DeleteProducer(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM producers WHERE id = $1", id)
	return err
}

func (r *Repository) GetProducerCount(ctx context.Context, userID int64) (int32, error) {
	var count int32
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM producers WHERE user_id = $1", userID).Scan(&count)
	return count, err
}
