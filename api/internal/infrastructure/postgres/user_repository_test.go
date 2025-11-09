package postgres_test

import (
	"context"
	"fmt"
	"hilo-api/internal/domain/do"
	"hilo-api/internal/infrastructure/postgres"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tdb := NewTestDB(t)
	defer tdb.Cleanup()

	repo := postgres.NewUserRepository(tdb.DB)
	ctx := context.Background()

	t.Run("create user successfully", func(t *testing.T) {
		user, err := do.NewUser("test@example.com", "password123", "testuser")
		require.NoError(t, err)

		err = repo.Create(ctx, user)
		assert.NoError(t, err)

		// Verify it was saved
		found, err := repo.FindByID(ctx, user.ID())
		require.NoError(t, err)
		assert.Equal(t, user.ID(), found.ID())
		assert.Equal(t, user.Email(), found.Email())
		assert.Equal(t, user.Username(), found.Username())
	})

	t.Run("cannot create duplicate email", func(t *testing.T) {
		user1, _ := do.NewUser("duplicate@example.com", "password123", "user1")
		user2, _ := do.NewUser("duplicate@example.com", "password456", "user2")

		err1 := repo.Create(ctx, user1)
		require.NoError(t, err1)

		err2 := repo.Create(ctx, user2)
		assert.Error(t, err2) // Should violate unique constraint
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tdb := NewTestDB(t)
	defer tdb.Cleanup()

	repo := postgres.NewUserRepository(tdb.DB)
	ctx := context.Background()

	t.Run("find existing user", func(t *testing.T) {
		user, _ := do.NewUser("find@example.com", "password123", "finduser")
		require.NoError(t, repo.Create(ctx, user))

		found, err := repo.FindByID(ctx, user.ID())

		require.NoError(t, err)
		assert.Equal(t, user.ID(), found.ID())
		assert.Equal(t, user.Email(), found.Email())
		assert.Equal(t, user.Username(), found.Username())
	})

	t.Run("user not found", func(t *testing.T) {
		nonExistentID := uuid.New()
		found, err := repo.FindByID(ctx, nonExistentID)

		assert.Error(t, err)
		assert.Nil(t, found)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestUserRepository_FindByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tdb := NewTestDB(t)
	defer tdb.Cleanup()

	repo := postgres.NewUserRepository(tdb.DB)
	ctx := context.Background()

	t.Run("find by email", func(t *testing.T) {
		user, _ := do.NewUser("email@example.com", "password123", "emailuser")
		require.NoError(t, repo.Create(ctx, user))

		found, err := repo.FindByEmail(ctx, "email@example.com")

		require.NoError(t, err)
		assert.Equal(t, user.ID(), found.ID())
		assert.Equal(t, user.Email(), found.Email())
	})

	t.Run("email not found", func(t *testing.T) {
		found, err := repo.FindByEmail(ctx, "nonexistent@example.com")

		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestUserRepository_FindAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tdb := NewTestDB(t)
	defer tdb.Cleanup()

	repo := postgres.NewUserRepository(tdb.DB)
	ctx := context.Background()

	// Create test users
	users := []*do.User{}
	for i := 0; i < 5; i++ {
		user, _ := do.NewUser(
			fmt.Sprintf("user%d@example.com", i),
			"password123",
			fmt.Sprintf("user%d", i),
		)
		require.NoError(t, repo.Create(ctx, user))
		users = append(users, user)
	}

	t.Run("find all with limit", func(t *testing.T) {
		found, err := repo.FindAll(ctx, 3, 0)

		require.NoError(t, err)
		assert.Len(t, found, 3)
	})

	t.Run("find all with offset", func(t *testing.T) {
		found, err := repo.FindAll(ctx, 3, 2)

		require.NoError(t, err)
		assert.Len(t, found, 3)
	})

	t.Run("find all returns empty when offset too large", func(t *testing.T) {
		found, err := repo.FindAll(ctx, 10, 100)

		require.NoError(t, err)
		assert.Empty(t, found)
	})
}

func TestUserRepository_Search(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tdb := NewTestDB(t)
	defer tdb.Cleanup()

	repo := postgres.NewUserRepository(tdb.DB)
	ctx := context.Background()

	// Create test users
	testUsers := []struct {
		email    string
		username string
	}{
		{"alice@example.com", "alice_wonder"},
		{"bob@example.com", "bob_builder"},
		{"charlie@example.com", "charlie_brown"},
		{"david@example.com", "alice_smith"},
	}

	for _, tu := range testUsers {
		user, _ := do.NewUser(tu.email, "password123", tu.username)
		require.NoError(t, repo.Create(ctx, user))
	}

	t.Run("search by username prefix", func(t *testing.T) {
		found, err := repo.Search(ctx, "alice", 10)

		require.NoError(t, err)
		assert.Len(t, found, 2)
	})

	t.Run("search case insensitive", func(t *testing.T) {
		found, err := repo.Search(ctx, "ALICE", 10)

		require.NoError(t, err)
		assert.Len(t, found, 2)
	})

	t.Run("search with limit", func(t *testing.T) {
		found, err := repo.Search(ctx, "alice", 1)

		require.NoError(t, err)
		assert.Len(t, found, 1)
	})

	t.Run("search no results", func(t *testing.T) {
		found, err := repo.Search(ctx, "nonexistent", 10)

		require.NoError(t, err)
		assert.Empty(t, found)
	})
}
