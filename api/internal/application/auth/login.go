package auth

import (
	"context"
	usecase "hilo-api/internal/application"
	"hilo-api/internal/domain/do"
	"hilo-api/internal/domain/repository"
)

// LoginUseCase handles user authentication
type LoginUseCase struct {
	userRepo repository.UserRepository
}

// NewLoginUseCase creates a new login use case
func NewLoginUseCase(userRepo repository.UserRepository) *LoginUseCase {
	return &LoginUseCase{
		userRepo: userRepo,
	}
}

// Execute authenticates a user
func (uc *LoginUseCase) Execute(ctx context.Context, email, password string) (*do.User, error) {
	// Find user
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, usecase.ErrInvalidCredentials
	}

	// Verify password (business rule in domain)
	if err := user.VerifyPassword(password); err != nil {
		return nil, usecase.ErrInvalidCredentials
	}

	return user, nil
}
