##Database Migrations Tool

Usage:
  go run cmd/migrate/main.go <command>

Commands:
  up       Apply all pending migrations (create tables)
  down     Rollback last migration  
  reset    Rollback all migrations
  status   Show migration status
  version  Show current version

Examples:
  go run cmd/migrate/main.go up     # Create all tables
  go run cmd/migrate/main.go down   # Remove last table
  go run cmd/migrate/main.go status # Show migration status

#Install goose
https://pressly.github.io/goose/installation