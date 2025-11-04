package auth

import (
	"context"
	"hilo-api/internal/domain"
)

// LoginUseCase handles user authentication
type LoginUseCase struct {
	userRepo domain.UserRepository
}

// NewLoginUseCase creates a new login use case
func NewLoginUseCase(userRepo domain.UserRepository) *LoginUseCase {
	return &LoginUseCase{
		userRepo: userRepo,
	}
}

// Execute authenticates a user
func (uc *LoginUseCase) Execute(ctx context.Context, email, password string) (*domain.User, error) {
	// Find user
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password (business rule in domain)
	if err := user.VerifyPassword(password); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
