package postgres_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

const migrationDir = "deployments/migrations/postgresql"

var (
	testDB *sqlx.DB
)

// TestMain sets up test database container
func TestMain(m *testing.M) {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	ctx := context.Background()

	// Start PostgreSQL container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
				"POSTGRES_DB":       "hilo_test",
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections"),
		},
	})
	if err != nil {
		logger.Fatal("failed to start postgres container", zap.Error(err))
	}

	defer container.Terminate(ctx)

	time.Sleep(2 * time.Second)

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	// Connect to database
	dsn := fmt.Sprintf("host=%s port=%d user=postgres password=postgres dbname=hilo_test sslmode=disable",
		host, port.Int())

	testDB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Fatal("failed to connect to test database", zap.Error(err))
	}

	// Run migrations
	if err := runMigrations(testDB, "hilo_test"); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	time.Sleep(1 * time.Second)

	// Run tests
	code := m.Run()

	// Cleanup
	testDB.Close()
	os.Exit(code)
}

// runMigrations applies database migrations
func runMigrations(db *sqlx.DB, dbName string) error {
	// Read SQL files directly instead of using golang-migrate
	// This is more reliable on Windows

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Construct absolute migration path
	migrationPath := filepath.Join(wd, "..", "..", "..", migrationDir)
	absPath, err := filepath.Abs(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	zap.L().Info("running migrations", zap.String("path", absPath))

	// Execute migration files in order
	migrationFiles := []string{
		"20251104031423_extension.up.sql",
		"20251104031532_database.up.sql",
	}

	for _, filename := range migrationFiles {
		filePath := filepath.Join(absPath, filename)

		sqlBytes, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		zap.L().Info("executing migration", zap.String("file", filename))

		if _, err := db.Exec(string(sqlBytes)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}
	}

	zap.L().Info("migrations completed successfully")
	return nil
}

// TestDB wraps sqlx.DB for test helpers
type TestDB struct {
	*sqlx.DB
	t *testing.T
}

// NewTestDB returns the shared test database
func NewTestDB(t *testing.T) *TestDB {
	t.Helper()
	return &TestDB{DB: testDB, t: t}
}

// Cleanup truncates all tables
func (tdb *TestDB) Cleanup() {
	tdb.t.Helper()

	_, err := tdb.Exec("TRUNCATE users, messages CASCADE")
	if err != nil {
		tdb.t.Fatalf("failed to cleanup database: %v", err)
	}
}
