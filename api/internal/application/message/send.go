package message

import (
	"context"
	"hilo-api/internal/domain"
	"hilo-api/internal/domain/repository"

	"github.com/google/uuid"
)

// SendMessageUseCase handles sending messages
type SendMessageUseCase struct {
	messageRepo repository.MessageRepository
	userRepo    repository.UserRepository
}

// NewSendMessageUseCase creates a new send message use case
func NewSendMessageUseCase(messageRepo repository.MessageRepository, userRepo repository.UserRepository) *SendMessageUseCase {
	return &SendMessageUseCase{
		messageRepo: messageRepo,
		userRepo:    userRepo,
	}
}

// Execute sends a message from sender to receiver
func (uc *SendMessageUseCase) Execute(ctx context.Context, senderID, receiverID uuid.UUID, content string) (*domain.Message, error) {
	// Verify receiver exists
	_, err := uc.userRepo.FindByID(ctx, receiverID)
	if err != nil {
		return nil, ErrReceiverNotFound
	}

	// Create message with business rules
	msg, err := domain.NewMessage(senderID, receiverID, content)
	if err != nil {
		return nil, err
	}

	// Persist
	if err := uc.messageRepo.Create(ctx, msg); err != nil {
		return nil, err
	}

	return msg, nil
}
