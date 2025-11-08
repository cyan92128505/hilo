package repository

import (
	"context"
	"hilo-api/internal/domain"

	"github.com/google/uuid"
)

// UserRepository defines user persistence operations
type UserRepository interface {
	// Create saves a new user
	Create(ctx context.Context, user *domain.User) error

	// FindByID retrieves user by ID
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)

	// FindByEmail retrieves user by email
	FindByEmail(ctx context.Context, email string) (*domain.User, error)

	// FindAll retrieves all users (for chat list)
	FindAll(ctx context.Context, limit, offset int) ([]*domain.User, error)

	// Search users by username
	Search(ctx context.Context, query string, limit int) ([]*domain.User, error)
}
