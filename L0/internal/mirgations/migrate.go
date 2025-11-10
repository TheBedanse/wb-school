package mirgations

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose"
)

func AutoMigrate(dbPassword string) error {
	connStr := fmt.Sprintf("postgres://L0User:%s@localhost:5432/L0?sslmode=disable", dbPassword)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer db.Close()

	return goose.Up(db, "schema")
}
