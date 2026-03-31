package repository

import (
	"context"
	"github.com/sig-agro/services/user-service/internal/entity"
)

type (
	UserRepository interface {
		Create(ctx context.Context, user *entity.User) error
		GetByID(ctx context.Context, id int64) (*entity.User, error)
		GetByEmail(ctx context.Context, email string) (*entity.User, error)
		ListAll(ctx context.Context) ([]entity.User, error)
		Update(ctx context.Context, user *entity.User) error
		Delete(ctx context.Context, id int64) error
		AddRole(ctx context.Context, userID int64, role string) error
		GetRoles(ctx context.Context, userID int64) ([]string, error)
	}
	CacheRepository interface {
		GetUser(id int64) (bool, *entity.User)
		SetUser(id int64, u *entity.User)
	}
)
