package routes

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"test2-api/handlers"
)

// corsMiddleware adds CORS headers to allow frontend access
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SetupRoutes configures all API routes
func SetupRoutes(db *sql.DB) http.Handler {
	router := mux.NewRouter()

	// Initialize handler
	feedbackHandler := handlers.NewFeedbackHandler(db)

	// API routes first (before static file server)
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/feedback", feedbackHandler.CreateFeedback).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/feedback", feedbackHandler.GetAllFeedback).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/feedback/{id}", feedbackHandler.GetFeedbackByID).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/feedback/{id}", feedbackHandler.DeleteFeedback).Methods("DELETE", "OPTIONS")

	// Health check endpoint
	apiRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API is running"))
	}).Methods("GET", "OPTIONS")

	// Serve frontend static files
	execPath, _ := os.Executable()
	execDir := filepath.Dir(execPath)
	
	// Try multiple possible paths for the frontend directory
	frontendPaths := []string{
		filepath.Join(execDir, "frontend"),
		filepath.Join(".", "frontend"),
	}

	var frontendDir string
	for _, path := range frontendPaths {
		if _, err := os.Stat(path); err == nil {
			frontendDir = path
			break
		}
	}

	if frontendDir != "" {
		// Serve static files at root - this catches all non-API routes
		router.PathPrefix("/").Handler(http.FileServer(http.Dir(frontendDir)))
	}

	// Apply CORS middleware
	return corsMiddleware(router)
}