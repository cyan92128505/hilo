package postgres

import (
	"context"
	"fmt"
	"time"

	"hilo-api/pkg/config"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func NewPostgresDB(logger *zap.Logger, opt config.Postgres) (*sqlx.DB, func(), error) {
	dsn := buildDSN(opt)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("postgres connect failed: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(opt.PostgresMaxOpenConns)
	db.SetMaxIdleConns(opt.PostgresMaxIdleConns)
	db.SetConnMaxLifetime(opt.PostgresConnMaxLifetime)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("postgres ping failed: %w", err)
	}

	logger.Info("postgres connected",
		zap.String("host", opt.PostgresHost),
		zap.String("database", opt.PostgresDatabase))

	cleanup := func() {
		if err := db.Close(); err != nil {
			logger.Error("postgres cleanup failed", zap.Error(err))
		}
	}

	return db, cleanup, nil
}

func buildDSN(opt config.Postgres) string {
	if opt.PostgresURL != "" {
		return opt.PostgresURL
	}
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s TimeZone=UTC",
		opt.PostgresUsername,
		opt.PostgresPassword,
		opt.PostgresHost,
		opt.PostgresPort,
		opt.PostgresDatabase,
		opt.PostgresSSLMode,
	)
}
