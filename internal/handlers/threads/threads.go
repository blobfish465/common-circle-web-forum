package threads

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"github.com/go-chi/chi/v5"
	"log"

	"github.com/blobfish465/common-circle-web-forum/internal/api"
	"github.com/blobfish465/common-circle-web-forum/internal/dataaccess/threads"
	"github.com/blobfish465/common-circle-web-forum/internal/database"
	"github.com/blobfish465/common-circle-web-forum/internal/models"
	"github.com/pkg/errors"
)

const (
	ListThreads                 = "threads.HandleList"
	SuccessfulListThreadsMessage = "Successfully listed threads"
	ErrRetrieveDatabase          = "Failed to retrieve database in %s"
	ErrRetrieveThreads           = "Failed to retrieve threads in %s"
	ErrEncodeView                = "Failed to encode threads in %s"
)

// ListThreads
func HandleListThreads(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	log.Println("Handling /threads request...")

	// Step 1: Get database connection
	db, err := database.GetDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return nil, errors.Wrap(err, "failed to connect to the database")
	}
	defer db.Close()

	// Step 2: Fetch threads
	threadsList, err := threads.List(db) 
	if err != nil {
		log.Println("Error fetching threads:", err)
		return nil, errors.Wrap(err, "failed to retrieve threads")
	}

	// Step 3: Encode threads to JSON
	data, err := json.Marshal(threadsList)
	if err != nil {
		log.Println("Error encoding threads to JSON:", err)
		return nil, errors.Wrap(err, "failed to encode threads")
	}

	// Step 4: Return API response
	response := &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{"Threads retrieved successfully"},
	}

	return response, nil
}

// Get a thread by ID
func HandleGetThreadByID(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	threadIDStr := chi.URLParam(r, "id")
	threadID, err := strconv.Atoi(threadIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid thread ID: %w", err)
	}

	db, err := database.GetDB()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve database: %w", err)
	}
	// Get thread from database 
	thread, err := threads.GetThreadByID(db, threadID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve thread with ID %d: %w", threadID, err)
	}

	data, err := json.Marshal(thread)
	if err != nil {
		return nil, fmt.Errorf("failed to encode thread data: %w", err)
	}

	return &api.Response{
		Payload: api.Payload{Data: data},
		Messages: []string{fmt.Sprintf("Thread retrieved successfully with ID %d", threadID)},
	}, nil
}

// List all of the threads created by a user
func HandleListThreadsByUser(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	userIDStr := chi.URLParam(r, "userId")
	if userIDStr == "" {
        return nil, fmt.Errorf("user ID is invalid or missing")
    }
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	db, err := database.GetDB()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve database: %w", err)
	}

	threads, err := threads.ListByUserID(db, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve threads for user %d: %w", userID, err)
	}

	data, err := json.Marshal(threads)
	if err != nil {
		return nil, fmt.Errorf("failed to encode threads: %w", err)
	}

	return &api.Response{
		Payload: api.Payload{Data: data},
		Messages: []string{fmt.Sprintf("Threads retrieved successfully for user %d", userID)},
	}, nil
}


// Handles creation of threads
func HandleCreateThreads(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	var thread models.Thread
	if err := json.NewDecoder(r.Body).Decode(&thread); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode thread: %v", err), http.StatusBadRequest)
		return nil, err
	}

	db, err := database.GetDB()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve database: %v", err), http.StatusInternalServerError)
		return nil, err
	}

	id, err := threads.Create(db, &thread)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create thread: %v", err), http.StatusInternalServerError)
		return nil, err
	}

	thread.ID = id
	data, _ := json.Marshal(thread)
	return &api.Response{
		Payload: api.Payload{Data: data},
		Messages: []string{"Thread created successfully"},
	}, nil
}

// Handles update of threads
func HandleUpdateThreads(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	// Extract the thread ID using path parameters
	threadIDStr := chi.URLParam(r, "id")
	threadID, err := strconv.Atoi(threadIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid thread ID: %w", err)
	}

	// Decode request body to get updated thread details
	var thread models.Thread
	if err := json.NewDecoder(r.Body).Decode(&thread); err != nil {
		return nil, fmt.Errorf("failed to decode thread: %w", err)
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

	originalThread, err := threads.GetThreadByID(db, threadID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch thread: %w", err)
	}

	// Ensure the user owns the thread
	if originalThread.UserID != userID {
		http.Error(w, "You are not authorized to update this thread", http.StatusForbidden)
		return nil, nil
	}

	// Set the thread ID for the updated thread
	thread.ID = threadID


	// Call update function in dataaccess thread
	err = threads.Update(db, &thread)
	if err != nil {
		return nil, fmt.Errorf("failed to update thread: %w", err)
	}

	return &api.Response{Messages: []string{"Thread updated successfully"}}, nil
}

// Handles deletion of threads
func HandleDeleteThreads(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	threadIDStr := chi.URLParam(r, "id")
	// Convert threadID to an integer
	threadID, err := strconv.Atoi(threadIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid thread ID: %w", err)
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
	
	originalThread, err := threads.GetThreadByID(db, threadID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch thread: %w", err)
	}

	// Ensure the user owns the thread
	if originalThread.UserID != userID {
		http.Error(w, "You are not authorized to update this thread", http.StatusForbidden)
		return nil, nil
	}

	err = threads.Delete(db, threadID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete thread: %w", err)
	}

	return &api.Response{Messages: []string{"Thread deleted successfully"}}, nil
}

// Handles listing of threads by category, for filtering threads by category
func HandleListThreadsByCategory(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
    categoryIDStr := chi.URLParam(r, "id")
	// chi.URLParam always extract parameters for URL as strings so need convert to int
	categoryID, err := strconv.Atoi(categoryIDStr)
    if err != nil {
        return nil, fmt.Errorf("invalid category ID: %w", err)
    }

    db, err := database.GetDB()
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve database: %w", err)
    }

    threads, err := threads.ListByCategoryID(db, categoryID)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve threads for category %d: %w", categoryID, err)
    }

	// Return an empty array if there are no threads in the category
    if len(threads) == 0 {
        threads = []models.Thread{}
    }


	// Marshal threads(a slice of models.Thread) into JSON format
	data, err := json.Marshal(threads)
	if err != nil {
		return nil, fmt.Errorf("failed to encode threads: %w", err)
	}

    return &api.Response{
        Payload: api.Payload{
			Data: data,
		},
        Messages: []string{"Threads retrieved successfully"},
    }, nil
}
