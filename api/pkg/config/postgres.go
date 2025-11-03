package config

import "time"

// Postgres type
type Postgres struct {
	PostgresURL       string        `split_words:"true" default:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	PostgresTxTimeout time.Duration `split_words:"true" default:"5s"`
}
