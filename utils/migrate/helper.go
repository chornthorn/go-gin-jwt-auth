package migrate

import (
	"fmt"
	"os/exec"
	"strings"
)

func CreateNewMigration(name string) error {
	if name == "" {
		return fmt.Errorf("migration name cannot be empty")
	}

	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ToLower(name)

	cmd := exec.Command("migrate",
		"create",
		"-ext", "sql",
		"-dir", "migrations/postgres",
		"-seq",
		name,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create migration: %w\nOutput: %s", err, output)
	}

	fmt.Printf("Created new migration: %s\n", name)
	return nil
}
