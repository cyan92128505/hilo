package user

import (
	"context"
	"hilo-api/internal/domain"
	"hilo-api/internal/domain/repository"
)

// ListUsersUseCase handles listing all users
type ListUsersUseCase struct {
	userRepo repository.UserRepository
}

// NewListUsersUseCase creates a new list users use case
func NewListUsersUseCase(userRepo repository.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{
		userRepo: userRepo,
	}
}

// Execute retrieves all users with pagination
func (uc *ListUsersUseCase) Execute(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return uc.userRepo.FindAll(ctx, limit, offset)
}
