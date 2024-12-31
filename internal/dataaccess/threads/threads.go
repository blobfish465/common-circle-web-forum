package threads

/*type Thread struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty"`
	CategoryID int `json:"category_id"`
}*/

import (
	"fmt"
	"log"
	"database/sql"
	"github.com/blobfish465/common-circle-web-forum/internal/database"
	"github.com/blobfish465/common-circle-web-forum/internal/models"
)

// List retrieves all threads from the database, for home page where all threads are listed
func List(db *database.Database) ([]models.Thread, error) {
	log.Println("Executing query to fetch threads...")
	rows, err := db.DB.Query(`
		SELECT id, user_id, title, content, created_at, updated_at, category_id
		FROM threads
	`)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var threads []models.Thread
	for rows.Next() {
		var thread models.Thread
		err := rows.Scan(&thread.ID, &thread.UserID, &thread.Title, &thread.Content, &thread.CreatedAt, &thread.UpdatedAt, &thread.CategoryID)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		threads = append(threads, thread)
	}
	if rows.Err() != nil {
        log.Println("Row iteration error:", rows.Err())
        return nil, rows.Err()
    }
	log.Printf("Fetched threads: %+v\n", threads)
	return threads, nil
}

// Create thread functionality, inserts new thread into database
func Create(db *database.Database, thread *models.Thread) (int, error) {
	query := `
		INSERT INTO threads (user_id, title, content, created_at, updated_at, category_id)
		VALUES ($1, $2, $3, NOW(), NOW(), $4) RETURNING id
	`
	var id int
	err := db.DB.QueryRow(query, thread.UserID, thread.Title, thread.Content, thread.CategoryID).Scan(&id)
	return id, err
}

// Update thread functionality, update existing thread in database
func Update(db *database.Database, thread *models.Thread) error {
	query := `
		UPDATE threads
		SET title = $1, content = $2, category_id = $3, updated_at = NOW()
		WHERE id = $4
	`
	_, err := db.DB.Exec(query, thread.Title, thread.Content, thread.CategoryID, thread.ID)
	return err
}

// Delete thread functionality, delete existing thread in database by its ID
func Delete(db *database.Database, threadID int) error {
	query := `
		DELETE FROM threads
		WHERE id = $1 
	`
	_, err := db.DB.Exec(query, threadID)
	return err
}

// Get thread from the database by its ID
func GetThreadByID(db *database.Database, id int) (*models.Thread, error) {
	var thread models.Thread
	query := `SELECT id, user_id, title, content, created_at, updated_at, category_id 
	FROM threads WHERE id = $1`
	row := db.DB.QueryRow(query, id)

	err := row.Scan(&thread.ID, &thread.UserID, &thread.Title, &thread.Content, &thread.CreatedAt, &thread.UpdatedAt, &thread.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("thread with ID %d not found", id)
		}
		return nil, fmt.Errorf("error retrieving thread with ID %d: %w", id, err)
	}

	return &thread, nil
}

// Get all the threads of a specific user
func ListByUserID(db *database.Database, userID int) ([]models.Thread, error) {
	rows, err := db.DB.Query(`
		SELECT id, user_id, title, content, created_at, updated_at, category_id
		FROM threads
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve threads for user %d: %w", userID, err)
	}
	defer rows.Close()

	var threads []models.Thread
	for rows.Next() {
		var thread models.Thread
		err := rows.Scan(&thread.ID, &thread.UserID, &thread.Title, &thread.Content, &thread.CreatedAt, &thread.UpdatedAt, &thread.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan thread data: %w", err)
		}
		threads = append(threads, thread)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over threads: %w", err)
	}

	return threads, nil
}

// Get all the threads of a category based on category_id, for filtering by category
func ListByCategoryID(db *database.Database, categoryID int) ([]models.Thread, error) {
    rows, err := db.DB.Query(`
        SELECT id, user_id, title, content, created_at, updated_at, category_id
        FROM threads
        WHERE category_id = $1
    `, categoryID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var threads []models.Thread
    for rows.Next() {
        var thread models.Thread
        err := rows.Scan(&thread.ID, &thread.UserID, &thread.Title, &thread.Content, &thread.CreatedAt, &thread.UpdatedAt, &thread.CategoryID)
        if err != nil {
            return nil, err
        }
        threads = append(threads, thread)
    }
    return threads, nil
}


