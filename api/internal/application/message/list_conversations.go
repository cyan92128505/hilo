package message

import (
	"context"
	"hilo-api/internal/domain"

	"github.com/google/uuid"
)

// ListConversationsUseCase handles listing all conversations for a user
type ListConversationsUseCase struct {
	messageRepo domain.MessageRepository
}

// NewListConversationsUseCase creates a new list conversations use case
func NewListConversationsUseCase(messageRepo domain.MessageRepository) *ListConversationsUseCase {
	return &ListConversationsUseCase{
		messageRepo: messageRepo,
	}
}

// Execute retrieves all conversations for a user
func (uc *ListConversationsUseCase) Execute(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.ConversationPreview, error) {
	return uc.messageRepo.ListUserConversations(ctx, userID, limit, offset)
}
