package users

import (
	"fmt"
	"database/sql"
	"log"
	"github.com/blobfish465/common-circle-web-forum/internal/database"
	"github.com/blobfish465/common-circle-web-forum/internal/models"
)

func List(db *database.Database) ([]models.User, error) {
	rows, err := db.DB.Query("SELECT id, username, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
	/*users := []models.User{
		{
			ID:   1,
			Name: "CVWO",
		},
	}
	return users, nil*/
}

// to retrieve a user from the database by their ID. 
func GetUserByID(db *database.Database, id int) (*models.User, error) {
	query := `SELECT id, username, email, password_hash FROM users WHERE id = $1`

	row := db.DB.QueryRow(query, id)

	var user models.User

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	return &user, nil
}

// Create and add new user into the database
func Create(db *database.Database, user *models.User) error {
	// Use db.DB to access the actual *sql.DB instance
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`
	_, err := db.DB.Exec(query, user.Username, user.Email, user.PasswordHash)
	return err
}

/* to test by returning the new ID:
func CreateAndReturnID(db *sql.DB, user *models.User) (int, error) {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`
	var newID int
	err := db.QueryRow(query, user.Username, user.Email, user.PasswordHash).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}*/

// Delete user from database
func Delete(db *database.Database, userID int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := db.DB.Exec(query, userID)
	return err
}

// GetUserByUsername retrieves a user from the database by their username
func GetUserByUsername(db *database.Database, username string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash FROM users WHERE username = $1`
	row := db.DB.QueryRow(query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with username %s not found", username)
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}
	log.Printf("Fetched user: %+v\n", user) // Log the fetched user
	return &user, nil
}
