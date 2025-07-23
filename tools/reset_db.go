package tools

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"test-go-htmx/internal/database"
	"test-go-htmx/internal/models"
)

func ResetDBMain() {
	// Define the database directory and file
	dbDir := "data"
	dbPath := filepath.Join(dbDir, "app.db")

	// Remove the existing database file if it exists
	if _, err := os.Stat(dbPath); err == nil {
		fmt.Printf("Removing existing database file: %s\n", dbPath)
		if err := os.Remove(dbPath); err != nil {
			log.Fatalf("Failed to remove database file: %v", err)
		}
	}

	// Create the database directory if it doesn't exist
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Initialize the database
	fmt.Println("Initializing database...")
	if err := database.Initialize(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.DB.Close()

	// Create sample journal entries
	fmt.Println("Creating sample journal entries...")

	// Create a gratitude journal entry
	id1, err := models.CreateJournal(
		"Grateful for Nature",
		"Today I took a walk in the park and felt truly grateful for the beauty of nature. The trees were especially vibrant.",
		"gratitude",
	)
	if err != nil {
		log.Fatalf("Failed to create gratitude journal entry: %v", err)
	}
	fmt.Printf("Created gratitude journal entry with ID: %d\n", id1)

	// Create a frustrations journal entry
	id2, err := models.CreateJournal(
		"Difficult Day at Work",
		"Today was challenging with tight deadlines and technical issues. I felt frustrated when my code wouldn't compile correctly.",
		"frustrations",
	)
	if err != nil {
		log.Fatalf("Failed to create frustrations journal entry: %v", err)
	}
	fmt.Printf("Created frustrations journal entry with ID: %d\n", id2)

	// Create another gratitude journal entry
	id3, err := models.CreateJournal(
		"Family Dinner",
		"I'm grateful for the wonderful dinner with my family tonight. These moments of connection are precious.",
		"gratitude",
	)
	if err != nil {
		log.Fatalf("Failed to create gratitude journal entry: %v", err)
	}
	fmt.Printf("Created gratitude journal entry with ID: %d\n", id3)

	// Create another frustrations journal entry
	id4, err := models.CreateJournal(
		"Traffic Jam",
		"Was stuck in traffic for over an hour today. It was frustrating to waste so much time just sitting in my car.",
		"frustrations",
	)
	if err != nil {
		log.Fatalf("Failed to create frustrations journal entry: %v", err)
	}
	fmt.Printf("Created frustrations journal entry with ID: %d\n", id4)

	fmt.Println("Database reset and initialized with sample data successfully!")
}
