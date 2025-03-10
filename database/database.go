package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// DB represents the database connection and provides methods to interact with it
type DB struct {
	conn *sql.DB
}

// New creates a new database connection
func New(connectionString string) (*DB, error) {
	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{conn: conn}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}
