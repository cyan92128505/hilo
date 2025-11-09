package auth

import (
	"context"
	usecase "hilo-api/internal/application"
	"hilo-api/internal/domain/do"
	"hilo-api/internal/domain/repository"
)

// RegisterUseCase handles user registration
type RegisterUseCase struct {
	userRepo repository.UserRepository
}

// NewRegisterUseCase creates a new register use case
func NewRegisterUseCase(userRepo repository.UserRepository) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo: userRepo,
	}
}

// Execute registers a new user
func (uc *RegisterUseCase) Execute(ctx context.Context, email, password, username string) (*do.User, error) {
	// Check if user already exists
	existing, err := uc.userRepo.FindByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, usecase.ErrEmailAlreadyExists
	}

	// Create user with business rules
	user, err := do.NewUser(email, password, username)
	if err != nil {
		return nil, err
	}

	// Persist
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
