package message

import (
	"context"
	"hilo-api/internal/domain/repository"

	"github.com/google/uuid"
)

// MarkAsReadUseCase handles marking messages as read
type MarkAsReadUseCase struct {
	messageRepo repository.MessageRepository
}

// NewMarkAsReadUseCase creates a new mark as read use case
func NewMarkAsReadUseCase(messageRepo repository.MessageRepository) *MarkAsReadUseCase {
	return &MarkAsReadUseCase{
		messageRepo: messageRepo,
	}
}

// Execute marks a message as read
func (uc *MarkAsReadUseCase) Execute(ctx context.Context, messageID, readerID uuid.UUID) error {
	// Load message
	msg, err := uc.messageRepo.FindByID(ctx, messageID)
	if err != nil {
		return err
	}

	// Apply business rule
	if err := msg.MarkAsRead(readerID); err != nil {
		return err
	}

	// Persist
	return uc.messageRepo.UpdateReadAt(ctx, msg.ID(), *msg.ReadAt())
}
