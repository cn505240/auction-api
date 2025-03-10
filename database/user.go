package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetUserById retrieves a user from the database by their ID
func (db *DB) GetUserById(id uuid.UUID) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := User{}
	err := db.conn.QueryRow(query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}
