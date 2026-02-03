package database

import (
	"github.com/darksuei/suei-intelligence/internal/config"
	"github.com/darksuei/suei-intelligence/internal/domain"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/database/postgres"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/database/sqlite"
)

func Initialize(config *config.DatabaseConfig) {
	switch config.DatabaseType {
		case domain.DatabaseTypePostgres:
			postgres.Connect(config)
		case domain.DatabaseTypeSqlite:
			sqlite.Connect(config)
		default:
			sqlite.Connect(config) // Treat SQLite as Default
	}
}

func Migrate(config *config.DatabaseConfig) {
	switch config.DatabaseType {
		case domain.DatabaseTypePostgres:
			postgres.Migrate()
		case domain.DatabaseTypeSqlite:
			sqlite.Migrate()
		default:
			sqlite.Migrate() // Treat SQLite as Default
	}
}
