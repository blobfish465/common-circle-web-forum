package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/blobfish465/common-circle-web-forum/internal/api"
	"github.com/blobfish465/common-circle-web-forum/internal/dataaccess/users"
	"github.com/blobfish465/common-circle-web-forum/internal/database"
	"github.com/blobfish465/common-circle-web-forum/internal/models"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-chi/chi/v5"
)

const (
	ListUsers                 = "users.HandleList"
	CreateUser                = "users.HandleCreate"
	SuccessfulListUsersMessage = "Successfully listed users"
	SuccessfulCreateUserMessage = "User created successfully"
	ErrRetrieveDatabase        = "Failed to retrieve database in %s"
	ErrRetrieveUsers           = "Failed to retrieve users in %s"
	ErrEncodeView              = "Failed to encode users in %s"
	ErrDecodeRequestBody       = "Failed to decode request body in %s"
	ErrHashPassword            = "Failed to hash password in %s"
	ErrCreateUser              = "Failed to create user in %s"
)

// ListUsers
func HandleListUsers(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	db, err := database.GetDB()

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrRetrieveDatabase, ListUsers))
	}

	users, err := users.List(db)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrRetrieveUsers, ListUsers))
	}

	data, err := json.Marshal(users)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrEncodeView, ListUsers))
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{SuccessfulListUsersMessage},
	}, nil
}

// Handles getting a user by id
func HandleGetUserByID(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	// Extract user ID using path parameters
	userIDStr := chi.URLParam(r, "id")
	if userIDStr == "" {
		return nil, fmt.Errorf("missing user ID")
	}

	// Convert the user ID string to an integer
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve database")
	}

	user, err := users.GetUserByID(db, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	data, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user data: %w", err)
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{"User retrieved successfully"},
	}, nil
}

// expected structure of the request body
type UserCreateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Create a new user
func HandleCreateUsers(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	// Decode the incoming request body
	var req UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrDecodeRequestBody, CreateUser))
	}

	// Validate required fields
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("missing required fields: username, email, or password")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrHashPassword, CreateUser))
	}

	// Establish database connection
	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrRetrieveDatabase, CreateUser))
	}

	// Create the user in the database
	newUser := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	err = users.Create(db, &newUser)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrCreateUser, CreateUser))
	}

	return &api.Response{
		Messages: []string{SuccessfulCreateUserMessage},
	}, nil
}

// Delete user from database
func HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID using path parameters
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// database connection
	db, err := database.GetDB()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}

	// Call dataaccess function in dataaccess/user.go to delete the user
	err = users.Delete(db, userID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully"))
}