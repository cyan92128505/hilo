package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrCannotSendToSelf = errors.New("cannot send message to yourself")
	ErrEmptyContent     = errors.New("message content cannot be empty")
	ErrNotReceiver      = errors.New("only receiver can mark message as read")
)

// Message represents a chat message between two users
type Message struct {
	id         uuid.UUID
	senderID   uuid.UUID
	receiverID uuid.UUID
	content    string
	createdAt  time.Time
	readAt     *time.Time
}

// NewMessage creates a new message with business rules enforced
func NewMessage(senderID, receiverID uuid.UUID, content string) (*Message, error) {
	if senderID == receiverID {
		return nil, ErrCannotSendToSelf
	}

	if content == "" {
		return nil, ErrEmptyContent
	}

	return &Message{
		id:         uuid.New(),
		senderID:   senderID,
		receiverID: receiverID,
		content:    content,
		createdAt:  time.Now(),
	}, nil
}

// ReconstructMessage rebuilds message from database (no validation)
func ReconstructMessage(id, senderID, receiverID uuid.UUID, content string, createdAt time.Time, readAt *time.Time) *Message {
	return &Message{
		id:         id,
		senderID:   senderID,
		receiverID: receiverID,
		content:    content,
		createdAt:  createdAt,
		readAt:     readAt,
	}
}

// MarkAsRead marks the message as read by the receiver
func (m *Message) MarkAsRead(readerID uuid.UUID) error {
	if readerID != m.receiverID {
		return ErrNotReceiver
	}

	if m.readAt != nil {
		return nil // already read, idempotent
	}

	now := time.Now()
	m.readAt = &now
	return nil
}

// Getters
func (m *Message) ID() uuid.UUID         { return m.id }
func (m *Message) SenderID() uuid.UUID   { return m.senderID }
func (m *Message) ReceiverID() uuid.UUID { return m.receiverID }
func (m *Message) Content() string       { return m.content }
func (m *Message) CreatedAt() time.Time  { return m.createdAt }
func (m *Message) ReadAt() *time.Time    { return m.readAt }
func (m *Message) IsRead() bool          { return m.readAt != nil }
