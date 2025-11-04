package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMessage(t *testing.T) {
	senderID := uuid.New()
	receiverID := uuid.New()

	t.Run("create valid message", func(t *testing.T) {
		msg, err := NewMessage(senderID, receiverID, "Hello World")

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, msg.ID())
		assert.Equal(t, senderID, msg.SenderID())
		assert.Equal(t, receiverID, msg.ReceiverID())
		assert.Equal(t, "Hello World", msg.Content())
		assert.False(t, msg.IsRead())
		assert.Nil(t, msg.ReadAt())
		assert.WithinDuration(t, time.Now(), msg.CreatedAt(), time.Second)
	})

	t.Run("cannot send to self", func(t *testing.T) {
		msg, err := NewMessage(senderID, senderID, "Self message")

		assert.Error(t, err)
		assert.Equal(t, ErrCannotSendToSelf, err)
		assert.Nil(t, msg)
	})

	t.Run("content cannot be empty", func(t *testing.T) {
		msg, err := NewMessage(senderID, receiverID, "")

		assert.Error(t, err)
		assert.Equal(t, ErrEmptyContent, err)
		assert.Nil(t, msg)
	})
}

func TestMessage_MarkAsRead(t *testing.T) {
	senderID := uuid.New()
	receiverID := uuid.New()
	msg, _ := NewMessage(senderID, receiverID, "Test")

	t.Run("receiver can mark as read", func(t *testing.T) {
		err := msg.MarkAsRead(receiverID)

		require.NoError(t, err)
		assert.True(t, msg.IsRead())
		assert.NotNil(t, msg.ReadAt())
		assert.WithinDuration(t, time.Now(), *msg.ReadAt(), time.Second)
	})

	t.Run("sender cannot mark as read", func(t *testing.T) {
		msg, _ := NewMessage(senderID, receiverID, "Test")
		err := msg.MarkAsRead(senderID)

		assert.Error(t, err)
		assert.Equal(t, ErrNotReceiver, err)
		assert.False(t, msg.IsRead())
	})

	t.Run("stranger cannot mark as read", func(t *testing.T) {
		msg, _ := NewMessage(senderID, receiverID, "Test")
		strangerID := uuid.New()
		err := msg.MarkAsRead(strangerID)

		assert.Error(t, err)
		assert.Equal(t, ErrNotReceiver, err)
		assert.False(t, msg.IsRead())
	})

	t.Run("marking as read twice is idempotent", func(t *testing.T) {
		msg, _ := NewMessage(senderID, receiverID, "Test")

		err1 := msg.MarkAsRead(receiverID)
		require.NoError(t, err1)
		firstReadAt := *msg.ReadAt()

		time.Sleep(10 * time.Millisecond)

		err2 := msg.MarkAsRead(receiverID)
		require.NoError(t, err2)
		secondReadAt := *msg.ReadAt()

		// ReadAt should not change on second call
		assert.Equal(t, firstReadAt, secondReadAt)
	})
}

func TestReconstructMessage(t *testing.T) {
	id := uuid.New()
	senderID := uuid.New()
	receiverID := uuid.New()
	createdAt := time.Now().Add(-1 * time.Hour)
	readAt := time.Now()

	t.Run("reconstruct unread message", func(t *testing.T) {
		msg := ReconstructMessage(id, senderID, receiverID, "Content", createdAt, nil)

		assert.Equal(t, id, msg.ID())
		assert.Equal(t, senderID, msg.SenderID())
		assert.Equal(t, receiverID, msg.ReceiverID())
		assert.Equal(t, "Content", msg.Content())
		assert.Equal(t, createdAt, msg.CreatedAt())
		assert.False(t, msg.IsRead())
		assert.Nil(t, msg.ReadAt())
	})

	t.Run("reconstruct read message", func(t *testing.T) {
		msg := ReconstructMessage(id, senderID, receiverID, "Content", createdAt, &readAt)

		assert.True(t, msg.IsRead())
		assert.NotNil(t, msg.ReadAt())
		assert.Equal(t, readAt, *msg.ReadAt())
	})
}
