package postgres

import (
	"testing"

	"holi-api/pkg/config"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestBuildDSN(t *testing.T) {
	tests := []struct {
		name string
		opt  config.Postgres
		want string
	}{
		{
			name: "use url when provided",
			opt: config.Postgres{
				PostgresURL: "postgres://user:pass@localhost:5432/db",
			},
			want: "postgres://user:pass@localhost:5432/db",
		},
		{
			name: "build from parameters",
			opt: config.Postgres{
				PostgresUsername: "testuser",
				PostgresPassword: "testpass",
				PostgresHost:     "localhost",
				PostgresPort:     "5432",
				PostgresDatabase: "testdb",
				PostgresSSLMode:  "disable",
			},
			want: "user=testuser password=testpass host=localhost port=5432 dbname=testdb sslmode=disable TimeZone=UTC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildDSN(tt.opt)
			assert.Equal(t, tt.want, got)
		})
	}
}

// Integration test - requires real postgres
func TestNewPostgresDB_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	logger := zap.NewNop()
	opt := config.Postgres{
		PostgresHost:     "localhost",
		PostgresPort:     "5432",
		PostgresUsername: "postgres",
		PostgresPassword: "postgres",
		PostgresDatabase: "test",
		PostgresSSLMode:  "disable",
	}

	db, cleanup, err := NewPostgresDB(logger, opt)
	if err != nil {
		t.Skipf("postgres not available: %v", err)
	}
	defer cleanup()

	assert.NotNil(t, db)

	// Verify connection works
	err = db.Ping()
	assert.NoError(t, err)
}
