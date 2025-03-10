package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

type Migration struct {
	ID       int
	Filename string
	SQL      string
}

func RunMigrations(db *sql.DB) error {
	// Create migrations table if it doesn't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("error creating migrations table: %v", err)
	}

	// Get applied migrations
	appliedVersions := make(map[int]bool)
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("error querying migrations: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return fmt.Errorf("error scanning migration version: %v", err)
		}
		appliedVersions[version] = true
	}

	// Read migration files
	migrations, err := loadMigrationFiles("database/migrations")
	if err != nil {
		return fmt.Errorf("error loading migrations: %v", err)
	}

	// Sort migrations by ID
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].ID < migrations[j].ID
	})

	// Apply new migrations
	for _, migration := range migrations {
		if !appliedVersions[migration.ID] {
			tx, err := db.Begin()
			if err != nil {
				return fmt.Errorf("error starting transaction: %v", err)
			}

			if _, err := tx.Exec(migration.SQL); err != nil {
				tx.Rollback()
				return fmt.Errorf("error applying migration %d: %v", migration.ID, err)
			}

			if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", migration.ID); err != nil {
				tx.Rollback()
				return fmt.Errorf("error recording migration %d: %v", migration.ID, err)
			}

			if err := tx.Commit(); err != nil {
				return fmt.Errorf("error committing migration %d: %v", migration.ID, err)
			}

			fmt.Printf("Applied migration: %s\n", migration.Filename)
		}
	}

	return nil
}

func loadMigrationFiles(dir string) ([]Migration, error) {
	var migrations []Migration

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			content, err := os.ReadFile(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, err
			}

			var id int
			_, err = fmt.Sscanf(file.Name(), "%d_", &id)
			if err != nil {
				return nil, fmt.Errorf("invalid migration filename format: %s", file.Name())
			}

			migrations = append(migrations, Migration{
				ID:       id,
				Filename: file.Name(),
				SQL:      string(content),
			})
		}
	}

	return migrations, nil
}
