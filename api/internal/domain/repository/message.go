package repository

import (
	"context"
	"hilo-api/internal/domain/do"
	"time"

	"github.com/google/uuid"
)

// MessageRepository defines message persistence operations
type MessageRepository interface {
	// Create saves a new message
	Create(ctx context.Context, msg *do.Message) error

	// FindByID retrieves message by ID
	FindByID(ctx context.Context, id uuid.UUID) (*do.Message, error)

	// UpdateReadAt marks message as read
	UpdateReadAt(ctx context.Context, id uuid.UUID, readAt time.Time) error

	// ListConversation retrieves messages between two users
	ListConversation(ctx context.Context, userA, userB uuid.UUID, limit, offset int) ([]*do.Message, error)

	// ListUserConversations retrieves all conversations for a user
	// Returns the latest message from each conversation
	ListUserConversations(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*do.ConversationPreview, error)
}
