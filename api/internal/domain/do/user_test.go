package do

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	t.Run("create valid user", func(t *testing.T) {
		user, err := NewUser("test@example.com", "password123", "testuser")

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, user.ID())
		assert.Equal(t, "test@example.com", user.Email())
		assert.Equal(t, "testuser", user.Username())
		assert.NotEmpty(t, user.PasswordHash())
		assert.NotEqual(t, "password123", user.PasswordHash()) // should be hashed
		assert.WithinDuration(t, time.Now(), user.CreatedAt(), time.Second)
	})

	t.Run("email cannot be empty", func(t *testing.T) {
		user, err := NewUser("", "password123", "testuser")

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidEmail, err)
		assert.Nil(t, user)
	})

	t.Run("password too short", func(t *testing.T) {
		user, err := NewUser("test@example.com", "short", "testuser")

		assert.Error(t, err)
		assert.Equal(t, ErrWeakPassword, err)
		assert.Nil(t, user)
	})

	t.Run("password must be at least 8 characters", func(t *testing.T) {
		user, err := NewUser("test@example.com", "1234567", "testuser")

		assert.Error(t, err)
		assert.Equal(t, ErrWeakPassword, err)
		assert.Nil(t, user)
	})

	t.Run("username cannot be empty", func(t *testing.T) {
		user, err := NewUser("test@example.com", "password123", "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username")
		assert.Nil(t, user)
	})
}

func TestUser_VerifyPassword(t *testing.T) {
	user, err := NewUser("test@example.com", "correct_password", "testuser")
	require.NoError(t, err)

	t.Run("correct password", func(t *testing.T) {
		err := user.VerifyPassword("correct_password")
		assert.NoError(t, err)
	})

	t.Run("wrong password", func(t *testing.T) {
		err := user.VerifyPassword("wrong_password")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidCredentials, err)
	})

	t.Run("empty password", func(t *testing.T) {
		err := user.VerifyPassword("")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidCredentials, err)
	})
}

func TestReconstructUser(t *testing.T) {
	id := uuid.New()
	email := "test@example.com"
	passwordHash := "$2a$12$abcdefghijklmnopqrstuvwxyz"
	username := "testuser"
	createdAt := time.Now().Add(-24 * time.Hour)

	t.Run("reconstruct user from database", func(t *testing.T) {
		user := ReconstructUser(id, email, passwordHash, username, createdAt)

		assert.Equal(t, id, user.ID())
		assert.Equal(t, email, user.Email())
		assert.Equal(t, passwordHash, user.PasswordHash())
		assert.Equal(t, username, user.Username())
		assert.Equal(t, createdAt, user.CreatedAt())
	})
}

func TestPasswordHashing(t *testing.T) {
	t.Run("same password generates different hashes", func(t *testing.T) {
		user1, err1 := NewUser("user1@example.com", "same_password", "user1")
		user2, err2 := NewUser("user2@example.com", "same_password", "user2")

		require.NoError(t, err1)
		require.NoError(t, err2)

		// Hashes should be different due to salt
		assert.NotEqual(t, user1.PasswordHash(), user2.PasswordHash())

		// But both should verify correctly
		assert.NoError(t, user1.VerifyPassword("same_password"))
		assert.NoError(t, user2.VerifyPassword("same_password"))
	})
}
