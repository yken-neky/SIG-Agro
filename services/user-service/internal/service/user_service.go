package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sig-agro/services/user-service/internal/entity"
	"github.com/sig-agro/services/user-service/internal/repository"
)

type UserService struct {
	repo      repository.UserRepository
	cache     repository.CacheRepository
	jwtSecret string
}

func NewUserService(repo repository.UserRepository, cache repository.CacheRepository, jwtSecret string) *UserService {
	return &UserService{repo: repo, cache: cache, jwtSecret: jwtSecret}
}

func (s *UserService) Register(ctx context.Context, email, password, fullName, phone string) (*entity.User, error) {
	hash := sha256.Sum256([]byte(password))
	passwordHash := fmt.Sprintf("%x", hash)

	user := &entity.User{
		Email:        email,
		PasswordHash: passwordHash,
		FullName:     fullName,
		Phone:        phone,
		CreatedAt:    time.Now(),
	}

	err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Add default role
	s.repo.AddRole(ctx, user.ID, "user")

	return user, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*entity.User, string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}

	hash := sha256.Sum256([]byte(password))
	passwordHash := fmt.Sprintf("%x", hash)

	if passwordHash != user.PasswordHash {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Generate JWT
	expiresAt := time.Now().Add(time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"roles":   user.Roles,
		"exp":     expiresAt,
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, "", err
	}

	return user, tokenString, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	// Check cache first
	if found, user := s.cache.GetUser(id); found {
		return user, nil
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	s.cache.SetUser(id, user)
	return user, nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]entity.User, error) {
	return s.repo.ListAll(ctx)
}

func (s *UserService) UpdateUser(ctx context.Context, user *entity.User) error {
	return s.repo.Update(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
