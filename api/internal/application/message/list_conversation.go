package message

import (
	"context"
	"hilo-api/internal/domain"

	"github.com/google/uuid"
)

// ListConversationUseCase handles listing messages in a conversation
type ListConversationUseCase struct {
	messageRepo domain.MessageRepository
}

// NewListConversationUseCase creates a new list conversation use case
func NewListConversationUseCase(messageRepo domain.MessageRepository) *ListConversationUseCase {
	return &ListConversationUseCase{
		messageRepo: messageRepo,
	}
}

// Execute retrieves messages between two users
func (uc *ListConversationUseCase) Execute(ctx context.Context, userA, userB uuid.UUID, limit, offset int) ([]*domain.Message, error) {
	return uc.messageRepo.ListConversation(ctx, userA, userB, limit, offset)
}
