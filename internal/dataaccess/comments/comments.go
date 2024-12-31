package comments

/*type Comment struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty"`
	UserID int `json:"user_id"`
	ThreadID  int    `json:"thread_id"`
}*/

import (
	"fmt"
	"database/sql"
	"github.com/blobfish465/common-circle-web-forum/internal/database"
	"github.com/blobfish465/common-circle-web-forum/internal/models"
)

// List retrieves all comments from the database.
func List(db *database.Database) ([]models.Comment, error) {
	rows, err := db.DB.Query(`
		SELECT id, content, created_at, updated_at, user_id, thread_id 
		FROM comments
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt, &comment.UserID, &comment.ThreadID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// Get a comment from the database by its ID
func GetCommentByID(db *database.Database, id int) (*models.Comment, error) {
	var comment models.Comment
	query := `
		SELECT id, content, created_at, updated_at, user_id, thread_id 
		FROM comments 
		WHERE id = $1
	`

	row := db.DB.QueryRow(query, id)

	err := row.Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt, &comment.UserID, &comment.ThreadID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("comment with ID %d not found", id)
		}
		return nil, fmt.Errorf("error retrieving comment with ID %d: %w", id, err)
	}

	return &comment, nil
}

// retrieves all comments for a specific thread from the database
func ListCommentsByThread(db *database.Database, threadID int) ([]models.Comment, error) {
	query := `
		SELECT id, content, created_at, updated_at, user_id, thread_id 
		FROM comments 
		WHERE thread_id = $1
	`

	rows, err := db.DB.Query(query, threadID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query for thread ID %d: %w", threadID, err)
	}
	defer rows.Close()

	var commentsList []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt, &comment.UserID, &comment.ThreadID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment data: %w", err)
		}
		commentsList = append(commentsList, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over rows: %w", err)
	}

	// Return the list of comments
	return commentsList, nil
}

// retrieves all the comments made by a specific user
func ListCommentsByUserID(db *database.Database, userID int) ([]models.Comment, error) {
	query := `
		SELECT id, content, created_at, updated_at, user_id, thread_id
		FROM comments
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commentsList []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt, &comment.UserID, &comment.ThreadID)
		if err != nil {
			return nil, err
		}
		commentsList = append(commentsList, comment)
	}

	return commentsList, nil
}


// Create Comment functionality, inserts new comment into database
func Create(db *database.Database, comment *models.Comment) (int, error) {
	query := `
		INSERT INTO comments (content, created_at, updated_at, user_id, thread_id)
		VALUES ($1, NOW(), NOW(), $2, $3) RETURNING id
	`
	var id int
	err := db.DB.QueryRow(query, comment.Content, comment.UserID, comment.ThreadID).Scan(&id)
	return id, err
}

// Update Comment functionality, update existing comment in database
func Update(db *database.Database, comment *models.Comment) error {
	query := `
		UPDATE comments
		SET content = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := db.DB.Exec(query, comment.Content, comment.ID)
	return err
}

// Delete comment functionality, delete existing comment in database by ID
func Delete(db *database.Database, commentID int) error {
	query := `
		DELETE FROM comments
		WHERE id = $1
	`
	_, err := db.DB.Exec(query, commentID)
	return err
}