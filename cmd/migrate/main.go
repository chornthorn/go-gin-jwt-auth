package main

import (
	"jwt-auth-app/config"
	"jwt-auth-app/migrations"
	"jwt-auth-app/utils/migrate"
	"log"
)

func main() {
	// Load environment variables
	config.LoadConfig()

	// Parse commands
	opts := migrate.ParseCommands()

	// Create migration manager
	mgr, err := migrations.NewMigrationManager(migrations.MigrationConfig{
		Host:     config.AppConfig.Database.Host,
		Port:     config.AppConfig.Database.Port,
		User:     config.AppConfig.Database.User,
		Password: config.AppConfig.Database.Password,
		DBName:   config.AppConfig.Database.DBName,
		SSLMode:  config.AppConfig.Database.SSLMode,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Handle commands
	if err := migrate.HandleCommands(mgr, opts); err != nil {
		log.Fatal(err)
	}
}
