package postgres_test

import (
	"context"
	"hilo-api/internal/domain"
	"hilo-api/internal/infrastructure/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tdb := NewTestDB(t)
	defer tdb.Cleanup()

	messageRepo := postgres.NewMessageRepository(tdb.DB)
	userRepo := postgres.NewUserRepository(tdb.DB)
	ctx := context.Background()

	t.Run("create message successfully", func(t *testing.T) {
		// Create users first
		sender, _ := domain.NewUser("sender@example.com", "password123", "sender")
		receiver, _ := domain.NewUser("receiver@example.com", "password123", "receiver")
		require.NoError(t, userRepo.Create(ctx, sender))
		require.NoError(t, userRepo.Create(ctx, receiver))

		// Create message
		msg, err := domain.NewMessage(sender.ID(), receiver.ID(), "Hello World")
		require.NoError(t, err)

		err = messageRepo.Create(ctx, msg)
		assert.NoError(t, err)

		// Verify it was saved
		found, err := messageRepo.FindByID(ctx, msg.ID())
		require.NoError(t, err)
		assert.Equal(t, msg.ID(), found.ID())
		assert.Equal(t, msg.Content(), found.Content())
		assert.False(t, found.IsRead())
	})
}

func TestMessageRepository_FindByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tdb := NewTestDB(t)
	defer tdb.Cleanup()

	messageRepo := postgres.NewMessageRepository(tdb.DB)
	userRepo := postgres.NewUserRepository(tdb.DB)
	ctx := context.Background()

	t.Run("find existing message", func(t *testing.T) {
		// Create users
		sender, _ := domain.NewUser("sender@example.com", "password123", "sender")
		receiver, _ := domain.NewUser("receiver@example.com", "password123", "receiver")
		require.NoError(t, userRepo.Create(ctx, sender))
		require.NoError(t, userRepo.Create(ctx, receiver))

		// Create message
		msg, _ := domain.NewMessage(sender.ID(), receiver.ID(), "Test message")
		require.NoError(t, messageRepo.Create(ctx, msg))

		found, err := messageRepo.FindByID(ctx, msg.ID())

		require.NoError(t, err)
		assert.Equal(t, msg.ID(), found.ID())
		assert.Equal(t, msg.Content(), found.Content())
	})

	t.Run("message not found", func(t *testing.T) {
		sender, _ := domain.NewUser("s@example.com", "password123", "s")
		receiver, _ := domain.NewUser("r@example.com", "password123", "r")
		require.NoError(t, userRepo.Create(ctx, sender))
		require.NoError(t, userRepo.Create(ctx, receiver))

		nonExistentMsg, _ := domain.NewMessage(sender.ID(), receiver.ID(), "test")
		found, err := messageRepo.FindByID(ctx, nonExistentMsg.ID())

		assert.Error(t, err)
		assert.Nil(t, found)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestMessageRepository_UpdateReadAt(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tdb := NewTestDB(t)
	defer tdb.Cleanup()

	messageRepo := postgres.NewMessageRepository(tdb.DB)
	userRepo := postgres.NewUserRepository(tdb.DB)
	ctx := context.Background()

	t.Run("update read_at timestamp", func(t *testing.T) {
		// Create users
		sender, _ := domain.NewUser("sender@example.com", "password123", "sender")
		receiver, _ := domain.NewUser("receiver@example.com", "password123", "receiver")
		require.NoError(t, userRepo.Create(ctx, sender))
		require.NoError(t, userRepo.Create(ctx, receiver))

		// Create unread message
		msg, _ := domain.NewMessage(sender.ID(), receiver.ID(), "Unread message")
		require.NoError(t, messageRepo.Create(ctx, msg))

		// Mark as read
		readAt := time.Now()
		err := messageRepo.UpdateReadAt(ctx, msg.ID(), readAt)
		require.NoError(t, err)

		// Verify update
		found, err := messageRepo.FindByID(ctx, msg.ID())
		require.NoError(t, err)
		assert.True(t, found.IsRead())
		assert.NotNil(t, found.ReadAt())
		assert.WithinDuration(t, readAt, *found.ReadAt(), time.Second)
	})
}

func TestMessageRepository_ListConversation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tdb := NewTestDB(t)
	defer tdb.Cleanup()

	messageRepo := postgres.NewMessageRepository(tdb.DB)
	userRepo := postgres.NewUserRepository(tdb.DB)
	ctx := context.Background()

	// Create users
	userA, _ := domain.NewUser("usera@example.com", "password123", "userA")
	userB, _ := domain.NewUser("userb@example.com", "password123", "userB")
	userC, _ := domain.NewUser("userc@example.com", "password123", "userC")
	require.NoError(t, userRepo.Create(ctx, userA))
	require.NoError(t, userRepo.Create(ctx, userB))
	require.NoError(t, userRepo.Create(ctx, userC))

	// Create conversation between A and B
	messages := []struct {
		sender   domain.User
		receiver domain.User
		content  string
	}{
		{*userA, *userB, "A to B message 1"},
		{*userB, *userA, "B to A message 1"},
		{*userA, *userB, "A to B message 2"},
		{*userB, *userA, "B to A message 2"},
		{*userC, *userA, "C to A message"}, // Different conversation
	}

	for _, m := range messages {
		msg, _ := domain.NewMessage(m.sender.ID(), m.receiver.ID(), m.content)
		require.NoError(t, messageRepo.Create(ctx, msg))
	}

	t.Run("list conversation between two users", func(t *testing.T) {
		found, err := messageRepo.ListConversation(ctx, userA.ID(), userB.ID(), 10, 0)

		require.NoError(t, err)
		assert.Len(t, found, 4)

		// Should be ordered by created_at DESC
		assert.Equal(t, "B to A message 2", found[0].Content())
	})

	t.Run("list conversation is bidirectional", func(t *testing.T) {
		foundAB, err1 := messageRepo.ListConversation(ctx, userA.ID(), userB.ID(), 10, 0)
		foundBA, err2 := messageRepo.ListConversation(ctx, userB.ID(), userA.ID(), 10, 0)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.Equal(t, len(foundAB), len(foundBA))
	})

	t.Run("list conversation with limit", func(t *testing.T) {
		found, err := messageRepo.ListConversation(ctx, userA.ID(), userB.ID(), 2, 0)

		require.NoError(t, err)
		assert.Len(t, found, 2)
	})

	t.Run("list conversation with offset", func(t *testing.T) {
		found, err := messageRepo.ListConversation(ctx, userA.ID(), userB.ID(), 2, 2)

		require.NoError(t, err)
		assert.Len(t, found, 2)
	})

	t.Run("no messages between users returns empty", func(t *testing.T) {
		newUserA, _ := domain.NewUser("new1@example.com", "password123", "new1")
		newUserB, _ := domain.NewUser("new2@example.com", "password123", "new2")
		require.NoError(t, userRepo.Create(ctx, newUserA))
		require.NoError(t, userRepo.Create(ctx, newUserB))

		found, err := messageRepo.ListConversation(ctx, newUserA.ID(), newUserB.ID(), 10, 0)

		require.NoError(t, err)
		assert.Empty(t, found)
	})
}

func TestMessageRepository_ListUserConversations(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tdb := NewTestDB(t)
	defer tdb.Cleanup()

	userRepo := postgres.NewUserRepository(tdb.DB)
	messageRepo := postgres.NewMessageRepository(tdb.DB)
	ctx := context.Background()

	// Create test users
	userA, _ := domain.NewUser("usera@example.com", "password123", "userA")
	userB, _ := domain.NewUser("userb@example.com", "password123", "userB")
	userC, _ := domain.NewUser("userc@example.com", "password123", "userC")

	require.NoError(t, userRepo.Create(ctx, userA))
	require.NoError(t, userRepo.Create(ctx, userB))
	require.NoError(t, userRepo.Create(ctx, userC))

	// Create conversations
	// A <-> B (2 messages, 1 unread from B)
	msgAB1, _ := domain.NewMessage(userA.ID(), userB.ID(), "A to B 1")
	msgBA1, _ := domain.NewMessage(userB.ID(), userA.ID(), "B to A 1 (unread)")
	require.NoError(t, messageRepo.Create(ctx, msgAB1))
	time.Sleep(100 * time.Millisecond) // Ensure different timestamps
	require.NoError(t, messageRepo.Create(ctx, msgBA1))

	// A <-> C (1 message, read)
	time.Sleep(100 * time.Millisecond)
	msgCA1, _ := domain.NewMessage(userC.ID(), userA.ID(), "C to A 1")
	require.NoError(t, messageRepo.Create(ctx, msgCA1))
	require.NoError(t, messageRepo.UpdateReadAt(ctx, msgCA1.ID(), time.Now()))

	t.Run("list user conversations", func(t *testing.T) {
		previews, err := messageRepo.ListUserConversations(ctx, userA.ID(), 10, 0)

		require.NoError(t, err)
		assert.Len(t, previews, 2)

		// Should be ordered by latest message (most recent first)
		// C message was created last, so it should be first
		assert.Equal(t, "C to A 1", previews[0].LastMessage.Content())
		assert.Equal(t, userC.ID(), previews[0].OtherUser.ID())
		assert.Equal(t, 0, previews[0].UnreadCount)

		// B message was created second, so it should be second
		assert.Equal(t, "B to A 1 (unread)", previews[1].LastMessage.Content())
		assert.Equal(t, userB.ID(), previews[1].OtherUser.ID())
		assert.Equal(t, 1, previews[1].UnreadCount)
	})

	t.Run("list conversations with limit", func(t *testing.T) {
		previews, err := messageRepo.ListUserConversations(ctx, userA.ID(), 1, 0)

		require.NoError(t, err)
		assert.Len(t, previews, 1)
	})

	t.Run("user with no conversations returns empty", func(t *testing.T) {
		newUser, _ := domain.NewUser("new@example.com", "password123", "newuser")
		require.NoError(t, userRepo.Create(ctx, newUser))

		previews, err := messageRepo.ListUserConversations(ctx, newUser.ID(), 10, 0)

		require.NoError(t, err)
		assert.Empty(t, previews)
	})
}
