package database

import (
	"L0/internal/interfaces"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

var _ interfaces.Repository = (*Database)(nil)

type Database struct {
	Conn *pgx.Conn
}

func NewDB(dbPassword, hostName string) (interfaces.Repository, error) {
	connStr := fmt.Sprintf("postgres://L0User:%s@%s:5432/L0", dbPassword, hostName)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable ping DB: %w", err)
	}

	log.Println("Connected to database")
	return &Database{Conn: conn}, nil
}

func (db *Database) Close() {
	if db.Conn != nil {
		db.Conn.Close(context.Background())
		log.Println("Database disconnect")
	}
}
