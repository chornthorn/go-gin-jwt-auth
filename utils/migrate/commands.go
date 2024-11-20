package migrate

import (
	"flag"
	"fmt"
	"jwt-auth-app/migrations"
)

type CommandOptions struct {
	Up      bool
	Down    bool
	Version uint
	Force   int
	Create  string
}

type Command struct {
	Name        string
	Description string
	Action      func(*migrations.MigrationManager) error
}

func ParseCommands() CommandOptions {
	opts := CommandOptions{}

	flag.BoolVar(&opts.Up, "up", false, "Run migrations up")
	flag.BoolVar(&opts.Down, "down", false, "Run migrations down")
	flag.UintVar(&opts.Version, "version", 0, "Migrate to specific version")
	flag.IntVar(&opts.Force, "force", -1, "Force set specific version")
	flag.StringVar(&opts.Create, "create", "", "Create a new migration")

	flag.Parse()

	return opts
}

func GetCommands() map[string]Command {
	return map[string]Command{
		"up": {
			Name:        "up",
			Description: "Run all pending migrations",
			Action: func(m *migrations.MigrationManager) error {
				if err := m.Up(); err != nil {
					return fmt.Errorf("failed to run migrations up: %w", err)
				}
				fmt.Println("Successfully ran migrations up")
				return nil
			},
		},
		"down": {
			Name:        "down",
			Description: "Rollback all migrations",
			Action: func(m *migrations.MigrationManager) error {
				if err := m.Down(); err != nil {
					return fmt.Errorf("failed to run migrations down: %w", err)
				}
				fmt.Println("Successfully ran migrations down")
				return nil
			},
		},
		"version": {
			Name:        "version",
			Description: "Show current migration version",
			Action: func(m *migrations.MigrationManager) error {
				version, dirty, err := m.Version()
				if err != nil {
					return fmt.Errorf("failed to get version: %w", err)
				}
				fmt.Printf("Current migration version: %d (dirty: %v)\n", version, dirty)
				return nil
			},
		},
	}
}

func HandleCommands(mgr *migrations.MigrationManager, opts CommandOptions) error {
	if opts.Create != "" {
		return CreateNewMigration(opts.Create)
	}

	switch {
	case opts.Up:
		return GetCommands()["up"].Action(mgr)

	case opts.Down:
		return GetCommands()["down"].Action(mgr)

	case opts.Version > 0:
		return func(m *migrations.MigrationManager) error {
			if err := m.MigrateTo(opts.Version); err != nil {
				return fmt.Errorf("failed to migrate to version %d: %w", opts.Version, err)
			}
			fmt.Printf("Successfully migrated to version %d\n", opts.Version)
			return nil
		}(mgr)

	case opts.Force >= 0:
		return func(m *migrations.MigrationManager) error {
			if err := m.Force(opts.Force); err != nil {
				return fmt.Errorf("failed to force version %d: %w", opts.Force, err)
			}
			fmt.Printf("Successfully forced version to %d\n", opts.Force)
			return nil
		}(mgr)

	default:
		return GetCommands()["version"].Action(mgr)
	}
}
