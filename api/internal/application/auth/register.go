package auth

import (
	"context"
	"hilo-api/internal/domain"
)

// RegisterUseCase handles user registration
type RegisterUseCase struct {
	userRepo domain.UserRepository
}

// NewRegisterUseCase creates a new register use case
func NewRegisterUseCase(userRepo domain.UserRepository) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo: userRepo,
	}
}

// Execute registers a new user
func (uc *RegisterUseCase) Execute(ctx context.Context, email, password, username string) (*domain.User, error) {
	// Check if user already exists
	existing, err := uc.userRepo.FindByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Create user with business rules
	user, err := domain.NewUser(email, password, username)
	if err != nil {
		return nil, err
	}

	// Persist
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
