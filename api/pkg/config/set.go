package config

import (
	"fmt"
)

func NewCore(set Set) Core         { return set.Core }
func NewJWT(set Set) JWT           { return set.JWT }
func NewPostgres(set Set) Postgres { return set.Postgres }
func NewServer(set Set) Server     { return set.Server }

func NewSet() (Set, error) {
	set := Set{}

	// Load each configuration section
	configs := []interface{}{
		&set.Core,
		&set.JWT,
		&set.Postgres,
		&set.Server,
	}

	for _, cfg := range configs {
		if err := LoadFromEnv(cfg); err != nil {
			return set, fmt.Errorf("failed to load config for %T: %w", cfg, err)
		}
	}

	return set, nil
}

type Set struct {
	Core     Core
	JWT      JWT
	Postgres Postgres
	Server   Server
}
