package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/blobfish465/common-circle-web-forum/internal/handlers/users"
	"github.com/blobfish465/common-circle-web-forum/internal/handlers/threads"
	"github.com/blobfish465/common-circle-web-forum/internal/handlers/comments"
	"github.com/blobfish465/common-circle-web-forum/internal/handlers/categories"
	"github.com/blobfish465/common-circle-web-forum/internal/auth"
	"net/http"
	"fmt"
	"encoding/json"
)

// GetPublicRoutes returns a function to set up public routes
func GetPublicRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/login", auth.Login) 

		r.Post("/users", func(w http.ResponseWriter, req *http.Request) {
			response, err := users.HandleCreateUsers(w, req)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})

		r.Get("/threads", func(w http.ResponseWriter, req *http.Request) {
			response, err := threads.HandleListThreads(w, req)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})

		r.Get("/threads/{id}", func(w http.ResponseWriter, req *http.Request) {
			response, err := threads.HandleGetThreadByID(w, req)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})

		r.Get("/threads/{thread_id}/comments", func(w http.ResponseWriter, req *http.Request) {
			response, err := comments.HandleListCommentsByThread(w, req)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})

		r.Post("/threads", func(w http.ResponseWriter, req *http.Request) {
			response, err := threads.HandleCreateThreads(w, req)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})

		// Get all the categories
		r.Get("/categories", func(w http.ResponseWriter, req *http.Request) {
			response, err := categories.HandleListCategories(w, req)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})

		// Get all the threads of the specified category
		r.Get("/categories/{id}/threads", func(w http.ResponseWriter, req *http.Request) {
			response, err := threads.HandleListThreadsByCategory(w, req)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})
		
		// Get category details of a specific cateogry
		r.Get("/categories/{id}", func(w http.ResponseWriter, req *http.Request) {
			response, err := categories.HandleGetCategoryByID(w, req)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})
	}
}

// GetPrivateRoutes sets up private routes requiring authentication
func GetPrivateRoutes(r chi.Router) {

	r.Get("/users/{id}", func(w http.ResponseWriter, req *http.Request) {
			response, err := users.HandleGetUserByID(w, req)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})

	r.Get("/users/{userId}/threads", func(w http.ResponseWriter, req *http.Request) {
		response, err := threads.HandleListThreadsByUser(w, req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	r.Delete("/users/{id}", func(w http.ResponseWriter, req *http.Request) {
		users.HandleDeleteUser(w, req)
	})

	r.Put("/threads/{id}", func(w http.ResponseWriter, req *http.Request) {
		response, err := threads.HandleUpdateThreads(w, req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	r.Delete("/threads/{id}", func(w http.ResponseWriter, req *http.Request) {
		response, err := threads.HandleDeleteThreads(w, req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Get comments made by a specific user
	r.Get("/users/{userId}/comments", func(w http.ResponseWriter, req *http.Request) {
		response, err := comments.HandleListCommentsByUser(w, req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})


	r.Post("/comments", func(w http.ResponseWriter, req *http.Request) {
		response, err := comments.HandleCreateComments(w, req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	r.Put("/comments/{id}", func(w http.ResponseWriter, req *http.Request) {
		response, err := comments.HandleUpdateComments(w, req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	r.Delete("/comments/{id}", func(w http.ResponseWriter, req *http.Request) {
		response, err := comments.HandleDeleteComments(w, req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}
