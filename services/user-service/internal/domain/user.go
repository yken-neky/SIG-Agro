package domain

type User struct {
	ID        int64
	Email     string
	FullName  string
	Phone     string
	Roles     []string
	CreatedAt int64
}

type RegisterRequest struct {
	Email    string
	Password string
	FullName string
	Phone    string
}

type LoginRequest struct {
	Email    string
	Password string
}

type TokenClaims struct {
	UserID    int64
	Email     string
	Roles     []string
	ExpiresAt int64
}
