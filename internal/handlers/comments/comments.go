package comments

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"github.com/go-chi/chi/v5"

	"github.com/blobfish465/common-circle-web-forum/internal/api"
	"github.com/blobfish465/common-circle-web-forum/internal/dataaccess/comments"
	"github.com/blobfish465/common-circle-web-forum/internal/database"
	"github.com/blobfish465/common-circle-web-forum/internal/models"
	"github.com/pkg/errors"
)

const (
	ListComments                 = "comments.HandleList"
	SuccessfulListCommentsMessage = "Successfully listed comments"
	ErrRetrieveDatabase           = "Failed to retrieve database in %s"
	ErrRetrieveComments           = "Failed to retrieve comments in %s"
	ErrEncodeView                 = "Failed to encode comments in %s"
)

// ListComments
func HandleListComments(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrRetrieveDatabase, ListComments))
	}

	commentsList, err := comments.List(db)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrRetrieveComments, ListComments))
	}

	data, err := json.Marshal(commentsList)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrEncodeView, ListComments))
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{SuccessfulListCommentsMessage},
	}, nil
}

// Handles getting of comment by comment ID
func HandleGetCommentByID(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	commentIDStr := chi.URLParam(r, "id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid comment ID: %w", err)
	}

	db, err := database.GetDB()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve database: %w", err)
	}

	// Get the comment from the database by ID
	comment, err := comments.GetCommentByID(db, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve comment with ID %d: %w", commentID, err)
	}

	data, err := json.Marshal(comment)
	if err != nil {
		return nil, fmt.Errorf("failed to encode comment data: %w", err)
	}

	return &api.Response{
		Payload: api.Payload{Data: data},
		Messages: []string{fmt.Sprintf("Comment retrieved successfully with ID %d", commentID)},
	}, nil
}


// Handles creation of comments 
func HandleCreateComments(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		return nil, fmt.Errorf("failed to decode comment: %w", err)
	}

	db, err := database.GetDB()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve database: %w", err)
	}

	id, err := comments.Create(db, &comment)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	comment.ID = id
	data, _ := json.Marshal(comment)
	return &api.Response{
		Payload: api.Payload{Data: data},
		Messages: []string{"Comment created successfully"},
	}, nil
}

// Handles listing of comments by thread ID
func HandleListCommentsByThread(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	// Extract the thread ID using path parameters
	threadIDStr := chi.URLParam(r, "thread_id")
	threadID, err := strconv.Atoi(threadIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid thread ID: %w", err)
	}

	// Establish a database connection
	db, err := database.GetDB()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve database: %w", err)
	}

	// Call the dataaccess function to get the comments related to the thread from the database
	commentsList, err := comments.ListCommentsByThread(db, threadID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve comments for thread ID %d: %w", threadID, err)
	}

	// Marshal the comments data into JSON
	data, err := json.Marshal(commentsList)
	if err != nil {
		return nil, fmt.Errorf("failed to encode comments data: %w", err)
	}

	// Return the response with the comments data
	return &api.Response{
		Payload: api.Payload{Data: data},
		Messages: []string{fmt.Sprintf("Comments retrieved successfully for thread ID %d", threadID)},
	}, nil
}

// Handles listing all comments made by a specific user
func HandleListCommentsByUser(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	// Extract the user ID from the URL
	userIDStr := chi.URLParam(r, "userId")
	if userIDStr == "" {
		return nil, fmt.Errorf("user ID is missing in the request")
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	db, err := database.GetDB()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve database: %w", err)
	}

	comments, err := comments.ListCommentsByUserID(db, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve comments for user %d: %w", userID, err)
	}

	// Encode the comments into JSON format
	data, err := json.Marshal(comments)
	if err != nil {
		return nil, fmt.Errorf("failed to encode comments: %w", err)
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{fmt.Sprintf("Comments retrieved successfully for user %d", userID)},
	}, nil
}

// Handles update of comments
func HandleUpdateComments(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	commentIDStr := chi.URLParam(r, "id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid comment ID: %w", err)
	}

	// Decode request body to get updated comment details
	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		return nil, fmt.Errorf("failed to decode comment: %w", err)
	}

	// Get user ID from request context (set by AuthMiddleware, user_id is a string)
	userIDStr, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "User ID is invalid", http.StatusUnauthorized)
		return nil, nil
	}
	userID, err := strconv.Atoi(userIDStr)


	db, err := database.GetDB()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve database: %w", err)
	}

	originalComment, err := comments.GetCommentByID(db, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comment: %w", err)
	}

	// Check if the user is the owner of the comment
	if originalComment.UserID != userID {
		http.Error(w, "You are not authorized to update this comment", http.StatusForbidden)
		return nil, nil
	}

	// Set the comment ID for the updated comment
	comment.ID = commentID

	err = comments.Update(db, &comment)
	if err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	return &api.Response{Messages: []string{"Comment updated successfully"}}, nil
}

// Handles deletion of comments
func HandleDeleteComments(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	// Extract comment ID from URL path
	commentIDStr := chi.URLParam(r, "id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid comment ID: %w", err)
	}

	// Get user ID from request context (set by AuthMiddleware, user_id is a string)
	userIDStr, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "User ID is invalid", http.StatusUnauthorized)
		return nil, nil
	}
	userID, err := strconv.Atoi(userIDStr)

	db, err := database.GetDB()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve database: %w", err)
	}

	originalComment, err := comments.GetCommentByID(db, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comment: %w", err)
	}

	// Check if the user is the owner of the comment
	if originalComment.UserID != userID {
		http.Error(w, "You are not authorized to delete this comment", http.StatusForbidden)
		return nil, nil
	}

	// Delete comment from the database
	err = comments.Delete(db, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete comment: %w", err)
	}

	return &api.Response{Messages: []string{"Comment deleted successfully"}}, nil
}