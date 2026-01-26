package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	Email        string     `db:"email" json:"email"`
	Phone        string     `db:"phone" json:"phone"`
	FullName     string     `db:"full_name" json:"full_name"`
	PasswordHash string     `db:"password_hash" json:"-"` // Never expose in JSON
	PINHash      *string    `db:"pin_hash" json:"-"`      // Never expose in JSON
	Status       UserStatus `db:"status" json:"status"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusBlocked   UserStatus = "blocked"
)

// IsActive checks if user is active
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// Validate validates user data
func (u *User) Validate() error {
	if u.Email == "" {
		return ErrInvalidEmail
	}
	if u.Phone == "" {
		return ErrInvalidPhone
	}
	if u.FullName == "" {
		return ErrInvalidFullName
	}
	return nil
}
