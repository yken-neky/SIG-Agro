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

func (r *Repository) CreateUser(ctx context.Context, email, passwordHash, fullName, phone string) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO users (email, password_hash, full_name, phone) VALUES ($1, $2, $3, $4) RETURNING id",
		email, passwordHash, fullName, phone,
	).Scan(&id)
	return id, err
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (int64, string, error) {
	var id int64
	var hash string
	err := r.db.QueryRowContext(ctx,
		"SELECT id, password_hash FROM users WHERE email = $1",
		email,
	).Scan(&id, &hash)
	return id, hash, err
}

func (r *Repository) GetUserByID(ctx context.Context, id int64) (string, string, string, error) {
	var email, fullName, phone string
	err := r.db.QueryRowContext(ctx,
		"SELECT email, full_name, phone FROM users WHERE id = $1",
		id,
	).Scan(&email, &fullName, &phone)
	return email, fullName, phone, err
}

func (r *Repository) AddUserRole(ctx context.Context, userID int64, role string) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO user_roles (user_id, role) VALUES ($1, $2) ON CONFLICT (user_id, role) DO NOTHING",
		userID, role,
	)
	return err
}

func (r *Repository) GetUserRoles(ctx context.Context, userID int64) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT role FROM user_roles WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (r *Repository) ListUsers(ctx context.Context, limit, offset int32) ([]map[string]interface{}, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, email, full_name, phone, created_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int64
		var email, fullName, phone string
		var createdAt int64
		if err := rows.Scan(&id, &email, &fullName, &phone, &createdAt); err != nil {
			return nil, err
		}
		users = append(users, map[string]interface{}{
			"id":         id,
			"email":      email,
			"full_name":  fullName,
			"phone":      phone,
			"created_at": createdAt,
		})
	}
	return users, rows.Err()
}
