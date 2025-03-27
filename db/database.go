package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
)

var DB *sql.DB

func InitDB() error {
	dbPath, err := getDatabasePath()
	if err != nil {
		return fmt.Errorf("failed to get database path: %w", err)
	}

	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory '%s': %w", dbDir, err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// log.Printf("\n*****Database connection established: %s*****", dbPath)

	DB = db

	return createTables()
}

func getDatabasePath() (string, error) {

	addDataDir, err := os.UserConfigDir()
	if err != nil {
		log.Println("Warning: Could not find user config directory. Using current directory.")
		return filepath.Join(".", "taskmanager.db"), nil
	}

	dbDir := filepath.Join(addDataDir, "taskmanager")
	dbPath := filepath.Join(dbDir, "taskmanager.db")

	return dbPath, nil
}

func createTables() error {

	createTasksTableSQL := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'pending', -- e.g., 'pending', 'done'
		created_at TEXT DEFAULT CURRENT_TIMESTAMP,
		due_date TEXT NULL -- Storing dates as ISO8601 strings (YYYY-MM-DD HH:MM:SS)
	);`

	createNotesTableSQL := `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
		created_at TEXT DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := DB.Exec(createTasksTableSQL); err != nil {
		return fmt.Errorf("failed to create tasks table: %w", err)
	}

	if _, err := DB.Exec(createNotesTableSQL); err != nil {
		return fmt.Errorf("failed to create notes table: %w", err)
	}

	// log.Println("Database tables checked/created successfully.")
	return nil

}

func CloseDB() {
	if DB != nil {
		DB.Close()
		// log.Println("\n*****Database connection closed.*****")
	}
}
