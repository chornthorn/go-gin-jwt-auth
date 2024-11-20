#!/bin/bash

# Function to display usage
show_usage() {
    echo "Usage: ./scripts/migrate.sh [command]"
    echo ""
    echo "Commands:"
    echo "  up              Run all pending migrations"
    echo "  down            Rollback all migrations"
    echo "  version         Show current migration version"
    echo "  to [version]    Migrate to specific version"
    echo "  force [version] Force set specific version"
    echo "  create [name]   Create a new migration"
    echo ""
    echo "Example:"
    echo "  ./scripts/migrate.sh up"
    echo "  ./scripts/migrate.sh create add_users_table"
}

# Make the script executable
chmod +x scripts/migrate.sh

# Handle commands
case "$1" in
    "up")
        go run cmd/migrate/main.go -up
        ;;
    "down")
        go run cmd/migrate/main.go -down
        ;;
    "version")
        go run cmd/migrate/main.go
        ;;
    "to")
        if [ -z "$2" ]; then
            echo "Error: Version number required"
            echo ""
            show_usage
            exit 1
        fi
        go run cmd/migrate/main.go -version "$2"
        ;;
    "force")
        if [ -z "$2" ]; then
            echo "Error: Version number required"
            echo ""
            show_usage
            exit 1
        fi
        go run cmd/migrate/main.go -force "$2"
        ;;
    "create")
        if [ -z "$2" ]; then
            echo "Error: Migration name required"
            echo ""
            show_usage
            exit 1
        fi
        go run cmd/migrate/main.go -create "$2"
        ;;
    *)
        show_usage
        exit 1
        ;;
esac