package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// DB is the global database connection
var DB *sql.DB

// Initialize sets up the database connection and runs migrations
func Initialize(dbPath string) error {
	var err error

	// Create the database directory if it doesn't exist
	dbDir := filepath.Dir(dbPath)
	if err = os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open the database connection
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Connected to database at %s", dbPath)

	// Run migrations
	if err = runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// runMigrations executes all SQL migration files in order
func runMigrations() error {
	migrationsDir := "internal/database/migrations"

	// Read migration files
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Sort and execute migrations
	migrations := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			migrations = append(migrations, filepath.Join(migrationsDir, entry.Name()))
		}
	}

	// Sort migrations by filename
	// (We're assuming migrations are named with a numeric prefix, e.g., 001_, 002_, etc.)

	for _, migrationPath := range migrations {
		log.Printf("Running migration: %s", migrationPath)

		// Read migration file
		migrationSQL, err := os.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", migrationPath, err)
		}

		// Execute migration
		_, err = DB.Exec(string(migrationSQL))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migrationPath, err)
		}

		log.Printf("Successfully applied migration: %s", migrationPath)
	}

	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return DB
}
