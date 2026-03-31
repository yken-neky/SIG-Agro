package entity

import "time"

type User struct {
	ID           int64      `cbor:"1,keyasint" json:"id"`
	Email        string     `cbor:"2,keyasint" json:"email"`
	PasswordHash string     `cbor:"3,keyasint" json:"password_hash"`
	FullName     string     `cbor:"4,keyasint" json:"full_name"`
	Phone        string     `cbor:"5,keyasint" json:"phone"`
	Roles        []string   `cbor:"6,keyasint" json:"roles"`
	CreatedAt    time.Time  `cbor:"7,keyasint" json:"created_at"`
	UpdatedAt    *time.Time `cbor:"8,keyasint" json:"updated_at"`
}
