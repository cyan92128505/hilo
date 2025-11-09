package message

import (
	"context"
	"hilo-api/internal/domain/do"
	"hilo-api/internal/domain/repository"

	"github.com/google/uuid"
)

// ListConversationUseCase handles listing messages in a conversation
type ListConversationUseCase struct {
	messageRepo repository.MessageRepository
}

// NewListConversationUseCase creates a new list conversation use case
func NewListConversationUseCase(messageRepo repository.MessageRepository) *ListConversationUseCase {
	return &ListConversationUseCase{
		messageRepo: messageRepo,
	}
}

// Execute retrieves messages between two users
func (uc *ListConversationUseCase) Execute(ctx context.Context, userA, userB uuid.UUID, limit, offset int) ([]*do.Message, error) {
	return uc.messageRepo.ListConversation(ctx, userA, userB, limit, offset)
}
