package database

import (
	"github.com/darksuei/suei-intelligence/internal/config"
	"github.com/darksuei/suei-intelligence/internal/domain"
	"github.com/darksuei/suei-intelligence/internal/domain/account"
	"github.com/darksuei/suei-intelligence/internal/domain/metadata"
	"github.com/darksuei/suei-intelligence/internal/domain/organization"
	"github.com/darksuei/suei-intelligence/internal/domain/project"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/database/postgres"
	postgresRepository "github.com/darksuei/suei-intelligence/internal/infrastructure/database/postgres/repositories"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/database/sqlite"
	sqliteRepository "github.com/darksuei/suei-intelligence/internal/infrastructure/database/sqlite/repositories"
	"gorm.io/gorm"
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

func GetDB(config *config.DatabaseConfig) *gorm.DB {
	switch config.DatabaseType {
		case domain.DatabaseTypePostgres:
			return postgres.DB
		case domain.DatabaseTypeSqlite:
			return sqlite.DB
		default:
			return sqlite.DB // Treat SQLite as Default
	}
}

func NewMetadataRepository(config *config.DatabaseConfig) metadata.MetadataRepository {
	db := GetDB(config)

	switch config.DatabaseType {
		case domain.DatabaseTypePostgres:
			return postgresRepository.NewMetadataRepository(db)
		case domain.DatabaseTypeSqlite:
			return sqliteRepository.NewMetadataRepository(db)
		default:
			return sqliteRepository.NewMetadataRepository(db) // Treat SQLite as Default
	}
}

func NewOrganizationRepository(config *config.DatabaseConfig) organization.OrganizationRepository {
	db := GetDB(config)

	switch config.DatabaseType {
		case domain.DatabaseTypePostgres:
			return postgresRepository.NewOrganizationRepository(db)
		case domain.DatabaseTypeSqlite:
			return sqliteRepository.NewOrganizationRepository(db)
		default:
			return sqliteRepository.NewOrganizationRepository(db) // Treat SQLite as Default
	}
}

func NewAccountRepository(config *config.DatabaseConfig) account.AccountRepository {
	db := GetDB(config)

	switch config.DatabaseType {
		case domain.DatabaseTypePostgres:
			return postgresRepository.NewAccountRepository(db)
		case domain.DatabaseTypeSqlite:
			return sqliteRepository.NewAccountRepository(db)
		default:
			return sqliteRepository.NewAccountRepository(db) // Treat SQLite as Default
	}
}

func NewProjectRepository(config *config.DatabaseConfig) project.ProjectRepository {
	db := GetDB(config)

	switch config.DatabaseType {
		case domain.DatabaseTypePostgres:
			return postgresRepository.NewProjectRepository(db)
		case domain.DatabaseTypeSqlite:
			return sqliteRepository.NewProjectRepository(db)
		default:
			return sqliteRepository.NewProjectRepository(db) // Treat SQLite as Default
	}
}