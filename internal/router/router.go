package router

import (
	"github.com/rs/cors"
	"github.com/go-chi/chi/v5"
	"github.com/blobfish465/common-circle-web-forum/internal/routes"
	"github.com/blobfish465/common-circle-web-forum/internal/middleware"
)

func Setup() chi.Router {
	// initialize router
	r := chi.NewRouter()

	// To allow cross origin requests, between localhost3000 frontend and localhost8000 backend
	// set up CORS(Cross-Origin Resource Sharing) middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://unique-brioche-acdf26.netlify.app/"}, // Frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // Allow cookies and other credentials
	})

	// Apply CORS middleware
	r.Use(corsMiddleware.Handler)

	setUpRoutes(r)
	return r
}

func setUpRoutes(r chi.Router) {
	// Public routes (no authentication needed)
	r.Group(routes.GetPublicRoutes())

	// Secured routes (requires authentication)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware) // Apply JWT authentication middleware
		routes.GetPrivateRoutes(r)      // Define secured routes here
	})
}
