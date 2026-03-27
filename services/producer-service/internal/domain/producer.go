package domain

type Producer struct {
	ID         int64
	UserID     int64
	Name       string
	DocumentID string
	Phone      string
	Email      string
	Address    string
	CreatedAt  int64
}
