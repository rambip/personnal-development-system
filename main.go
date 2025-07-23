package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"test-go-htmx/internal/database"
	"test-go-htmx/internal/handlers"
)

func main() {
	// Set up database
	dbDir := "data"
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	dbPath := filepath.Join(dbDir, "app.db")
	if err := database.Initialize(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Define the file server for static assets
	staticDir := "web/static"
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Define the routes
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/journals", handlers.JournalsHandler)
	http.HandleFunc("/values", handlers.ValuesHandler)
	http.HandleFunc("/values/children", handlers.ValuesHandler)
	http.HandleFunc("/values/parents", handlers.ValuesHandler)
	http.HandleFunc("/journals/type/", handlers.JournalsHandler)
	http.HandleFunc("/journals/", handlers.JournalDetailHandler)

	// Start the server
	port := ":8888"
	log.Printf("Starting server on http://localhost%s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
