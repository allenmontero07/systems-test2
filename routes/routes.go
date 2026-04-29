package routes

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const idKey contextKey = "id"

// extractIDFromPath extracts the ID from a path like /feedback/{id}
func extractIDFromPath(path, prefix string) string {
	// Remove the prefix (e.g., /api/feedback/)
	trimmed := strings.TrimPrefix(path, prefix)
	// Return the remaining part which is the ID
	return trimmed
}

// SetupRoutes configures all API routes using standard library http.ServeMux
func SetupRoutes(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// Initialize handler
	feedbackHandler := handlers.NewFeedbackHandler(db)

	// API routes
	mux.HandleFunc("POST /api/feedback", feedbackHandler.CreateFeedback)
	mux.HandleFunc("GET /api/feedback", feedbackHandler.GetAllFeedback)
	
	// For paths with parameters, we need to handle them manually
	mux.HandleFunc("GET /api/feedback/", func(w http.ResponseWriter, r *http.Request) {
		// Extract ID from path
		id := extractIDFromPath(r.URL.Path, "/api/feedback/")
		if id == "" {
			http.NotFound(w, r)
			return
		}
		// Set the ID in the request context for the handler to use
		r = r.WithContext(context.WithValue(r.Context(), idKey, id))
		feedbackHandler.GetFeedbackByID(w, r)
	})
	
	mux.HandleFunc("DELETE /api/feedback/", func(w http.ResponseWriter, r *http.Request) {
		// Extract ID from path
		id := extractIDFromPath(r.URL.Path, "/api/feedback/")
		if id == "" {
			http.NotFound(w, r)
			return
		}
		// Set the ID in the request context for the handler to use
		r = r.WithContext(context.WithValue(r.Context(), idKey, id))
		feedbackHandler.DeleteFeedback(w, r)
	})

	// Health check endpoint
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API is running"))
	})

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

	// Create a handler that serves static files for non-API routes
	staticHandler := http.FileServer(http.Dir(frontendDir))
	
	// Wrap static handler to only serve non-API requests
	wrappedStatic := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			http.NotFound(w, r)
			return
		}
		staticHandler.ServeHTTP(w, r)
	})
	
	mux.Handle("/", wrappedStatic)

	// Apply CORS middleware
	return corsMiddleware(mux)
}