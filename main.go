package main

import (
	"fmt"
	"log"
	"net/http"

	"test2-api/config"
	"test2-api/routes"
)

func main() {
	// Connect to the database
	config.ConnectDB()

	// Run migrations (optional, can be done manually)
	config.RunMigrations("migrations/")

	// Setup routes
	router := routes.SetupRoutes(config.DB)

	// Start server
	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}