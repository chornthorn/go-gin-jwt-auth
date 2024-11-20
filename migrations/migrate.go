package migrations

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type MigrationConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type MigrationManager struct {
	config  MigrationConfig
	migrate *migrate.Migrate
}

func NewMigrationManager(config MigrationConfig) (*MigrationManager, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
		config.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("error creating postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/postgres",
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating migration instance: %w", err)
	}

	return &MigrationManager{
		config:  config,
		migrate: m,
	}, nil
}

func (m *MigrationManager) Up() error {
	err := m.migrate.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error running migrations up: %w", err)
	}
	return nil
}

func (m *MigrationManager) Down() error {
	err := m.migrate.Down()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error running migrations down: %w", err)
	}
	return nil
}

func (m *MigrationManager) Version() (uint, bool, error) {
	return m.migrate.Version()
}

func (m *MigrationManager) MigrateTo(version uint) error {
	err := m.migrate.Migrate(version)
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error migrating to version %d: %w", version, err)
	}
	return nil
}

func (m *MigrationManager) Force(version int) error {
	err := m.migrate.Force(version)
	if err != nil {
		return fmt.Errorf("error forcing version %d: %w", version, err)
	}
	return nil
}
