package main

import (
	"L0/internal/config"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	command := os.Args[1]

	cfg := config.LoadConfig()
	if cfg.DBPassword == "" {
		log.Fatal("DB_PASSWORD not found")
	}

	connStr := fmt.Sprintf("postgres://L0User:%s@%s:5432/L0", cfg.DBPassword, cfg.HostName)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	migrationsDir := "schema"

	switch command {
	case "up":
		fmt.Println("Applying migrations...")
		if err := goose.Up(db, migrationsDir); err != nil {
			log.Fatal("Failed to apply migrations:", err)
		}

	case "down":
		fmt.Println("Rolling back last migration...")
		if err := goose.Down(db, migrationsDir); err != nil {
			log.Fatal("Failed to rollback migration:", err)
		}

	case "reset":
		fmt.Println("Resetting all migrations...")
		if err := goose.Reset(db, migrationsDir); err != nil {
			log.Fatal("Failed to reset migrations:", err)
		}

	case "status":
		fmt.Println("Migration status:")
		if err := goose.Status(db, migrationsDir); err != nil {
			log.Fatal("Failed to get status:", err)
		}

	case "version":
		version, err := goose.GetDBVersion(db)
		if err != nil {
			log.Fatal("Failed to get version:", err)
		}
		fmt.Printf("Current migration version: %d\n", version)

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
