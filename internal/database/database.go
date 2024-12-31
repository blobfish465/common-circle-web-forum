package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Database struct wraps a pointer to an sql.DB instance, which represents a pool of database connections.
// to extend database functionality later, such as adding helper methods for transactions or other operations.
type Database struct {
	DB *sql.DB
}

// Creates and returns a new database connection wrapped in a Database struct.
func GetDB() (*Database, error) {
	log.Println("Fetching DATABASE_URL from environment...")
	connStr := os.Getenv("DATABASE_URL")  // Fetch the database URL from environment variables
    if connStr == "" {
        log.Fatal("DATABASE_URL is not set in the environment variables")
    }

	log.Println("Connecting to database with connection string:", connStr)
    db, err := sql.Open("postgres", connStr)
    if err != nil {
		log.Println("Error opening database connection:", err)
        return nil, err
    }
    db.SetMaxOpenConns(1000)
    db.SetMaxIdleConns(500)
    db.SetConnMaxLifetime(30 * time.Minute)

    err = setupTables(db)
    if err != nil {
        return nil, err
    }

    return &Database{DB: db}, nil
}

func (db *Database) Close() {
	err := db.DB.Close()
	if err != nil {
		log.Println("Error closing database:", err)
	}
}

func insertPredefinedCategories(db *sql.DB) error {
	predefinedCategories := []string{
		"Technology",
		"Health",
		"Education",
		"Science",
		"Sports",
		"Travel",
		"Entertainment",
	}

	for _, category := range predefinedCategories {
		query := `
		INSERT INTO categories (name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING;
		`
		_, err := db.Exec(query, category)
		if err != nil {
			log.Printf("Error inserting category %s: %v\n", category, err)
			return err
		}
	}

	log.Println("Predefined categories inserted successfully.")
	return nil
}


// setupTables creates tables if they do not exist
func setupTables(db *sql.DB) error {
	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			password_hash TEXT NOT NULL
		);`,
		`
		CREATE TABLE IF NOT EXISTS categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE
		);`,
		`
		CREATE TABLE IF NOT EXISTS threads (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP,
			category_id INT REFERENCES categories(id) ON DELETE SET NULL
		);`,
		`
		CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			thread_id INT NOT NULL REFERENCES threads(id) ON DELETE CASCADE
		);`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}

	// Insert predefined categories
	err := insertPredefinedCategories(db)
	if err != nil {
		return err
	}

	log.Println("All tables are set up successfully.")
	return nil
}