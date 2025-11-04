package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

const (
	MinPasswordLength = 8
	BcryptCost        = 12
)

// User represents a user account
type User struct {
	id           uuid.UUID
	email        string
	passwordHash string
	username     string
	createdAt    time.Time
}

// NewUser creates a new user with hashed password
func NewUser(email, password, username string) (*User, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}

	if len(password) < MinPasswordLength {
		return nil, ErrWeakPassword
	}

	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return nil, err
	}

	return &User{
		id:           uuid.New(),
		email:        email,
		passwordHash: string(hash),
		username:     username,
		createdAt:    time.Now(),
	}, nil
}

// ReconstructUser rebuilds user from database (no validation)
func ReconstructUser(id uuid.UUID, email, passwordHash, username string, createdAt time.Time) *User {
	return &User{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		username:     username,
		createdAt:    createdAt,
	}
}

// VerifyPassword checks if the provided password matches
func (u *User) VerifyPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.passwordHash), []byte(password))
	if err != nil {
		return ErrInvalidCredentials
	}
	return nil
}

// Getters
func (u *User) ID() uuid.UUID        { return u.id }
func (u *User) Email() string        { return u.email }
func (u *User) PasswordHash() string { return u.passwordHash }
func (u *User) Username() string     { return u.username }
func (u *User) CreatedAt() time.Time { return u.createdAt }
