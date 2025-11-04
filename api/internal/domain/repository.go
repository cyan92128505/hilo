package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// UserRepository defines user persistence operations
type UserRepository interface {
	// Create saves a new user
	Create(ctx context.Context, user *User) error

	// FindByID retrieves user by ID
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)

	// FindByEmail retrieves user by email
	FindByEmail(ctx context.Context, email string) (*User, error)

	// FindAll retrieves all users (for chat list)
	FindAll(ctx context.Context, limit, offset int) ([]*User, error)

	// Search users by username
	Search(ctx context.Context, query string, limit int) ([]*User, error)
}

// MessageRepository defines message persistence operations
type MessageRepository interface {
	// Create saves a new message
	Create(ctx context.Context, msg *Message) error

	// FindByID retrieves message by ID
	FindByID(ctx context.Context, id uuid.UUID) (*Message, error)

	// UpdateReadAt marks message as read
	UpdateReadAt(ctx context.Context, id uuid.UUID, readAt time.Time) error

	// ListConversation retrieves messages between two users
	ListConversation(ctx context.Context, userA, userB uuid.UUID, limit, offset int) ([]*Message, error)

	// ListUserConversations retrieves all conversations for a user
	// Returns the latest message from each conversation
	ListUserConversations(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*ConversationPreview, error)
}

// ConversationPreview represents the latest message in a conversation
type ConversationPreview struct {
	OtherUser *User
	LastMessage *Message
	UnreadCount int
}
