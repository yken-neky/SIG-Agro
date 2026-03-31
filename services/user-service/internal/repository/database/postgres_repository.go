package database

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sig-agro/services/user-service/internal/entity"
	"github.com/sig-agro/services/user-service/internal/repository"
	"github.com/sig-agro/services/user-service/internal/repository/sqlc"
)

type PostgresRepository struct {
	q *sqlc.Queries
}

func NewPostgresRepository(db sqlc.DBTX) repository.UserRepository {
	return &PostgresRepository{q: sqlc.New(db)}
}

func pgTypeTextToString(t pgtype.Text) string {
	if t.Valid {
		return t.String
	}
	return ""
}

func stringToPgTypeText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

func (r *PostgresRepository) Create(ctx context.Context, user *entity.User) error {
	hash := sha256.Sum256([]byte(user.PasswordHash))
	passwordHash := fmt.Sprintf("%x", hash)

	params := sqlc.CreateUserParams{
		Email:        user.Email,
		PasswordHash: passwordHash,
		FullName: pgtype.Text{
			String: user.FullName,
			Valid:  true,
		},
		Phone: pgtype.Text{
			String: user.Phone,
			Valid:  true,
		},
	}

	created, err := r.q.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	user.ID = created.ID
	user.CreatedAt = created.CreatedAt.Time
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	u, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	roles, err := r.q.GetUserRoles(ctx, id)
	if err != nil {
		roles = []string{}
	}

	return &entity.User{
		ID:        u.ID,
		Email:     u.Email,
		FullName:  pgTypeTextToString(u.FullName),
		Phone:     pgTypeTextToString(u.Phone),
		Roles:     roles,
		CreatedAt: u.CreatedAt.Time,
	}, nil
}

func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	u, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	roles, err := r.q.GetUserRoles(ctx, u.ID)
	if err != nil {
		roles = []string{}
	}

	return &entity.User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		FullName:     pgTypeTextToString(u.FullName),
		Phone:        pgTypeTextToString(u.Phone),
		Roles:        roles,
		CreatedAt:    u.CreatedAt.Time,
	}, nil
}

func (r *PostgresRepository) ListAll(ctx context.Context) ([]entity.User, error) {
	users, err := r.q.ListUsers(ctx, sqlc.ListUsersParams{Limit: 100, Offset: 0})
	if err != nil {
		return nil, err
	}

	var result []entity.User
	for _, u := range users {
		result = append(result, entity.User{
			ID:        u.ID,
			Email:     u.Email,
			FullName:  pgTypeTextToString(u.FullName),
			Phone:     pgTypeTextToString(u.Phone),
			CreatedAt: u.CreatedAt.Time,
		})
	}
	return result, nil
}

func (r *PostgresRepository) Update(ctx context.Context, user *entity.User) error {
	return r.q.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:       user.ID,
		FullName: stringToPgTypeText(user.FullName),
		Phone:    stringToPgTypeText(user.Phone),
	})
}

func (r *PostgresRepository) Delete(ctx context.Context, id int64) error {
	// Assuming there's a delete query, but not in sql, so placeholder
	return fmt.Errorf("delete not implemented")
}

func (r *PostgresRepository) AddRole(ctx context.Context, userID int64, role string) error {
	return r.q.CreateUserRole(ctx, sqlc.CreateUserRoleParams{
		UserID: userID,
		Role:   role,
	})
}

func (r *PostgresRepository) GetRoles(ctx context.Context, userID int64) ([]string, error) {
	return r.q.GetUserRoles(ctx, userID)
}
