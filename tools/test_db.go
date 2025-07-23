package tools

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"pds/internal/database"
	"pds/internal/models"
)

func TestDBMain() {
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

	// Insert test journal entries
	fmt.Println("Creating test journal entries...")
	id1, err := models.CreateJournal("First Journal Entry", "This is the content of my first journal entry.", "gratitude")
	if err != nil {
		log.Fatalf("Failed to create journal entry: %v", err)
	}
	fmt.Printf("Created journal entry with ID: %d\n", id1)

	id2, err := models.CreateJournal("Second Journal Entry", "This is the content of my second journal entry.", "frustrations")
	if err != nil {
		log.Fatalf("Failed to create journal entry: %v", err)
	}
	fmt.Printf("Created journal entry with ID: %d\n", id2)

	// Retrieve all journal entries
	fmt.Println("\nRetrieving all journal entries...")
	journals, err := models.GetAllJournals()
	if err != nil {
		log.Fatalf("Failed to retrieve journal entries: %v", err)
	}

	// Display all journal entries
	fmt.Printf("Found %d journal entries:\n", len(journals))
	for _, journal := range journals {
		fmt.Printf("ID: %d\nTitle: %s\nContent: %s\nType: %s\nCreated: %v\nUpdated: %v\n\n",
			journal.ID, journal.Title, journal.Content, journal.JournalType, journal.CreatedAt, journal.UpdatedAt)
	}

	// Update a journal entry
	fmt.Println("Updating journal entry...")
	err = models.UpdateJournal(id1, "Updated First Entry", "This content has been updated.", "gratitude")
	if err != nil {
		log.Fatalf("Failed to update journal entry: %v", err)
	}

	// Retrieve and display the updated entry
	fmt.Println("Retrieving updated journal entry...")
	journal, err := models.GetJournal(id1)
	if err != nil {
		log.Fatalf("Failed to retrieve journal entry: %v", err)
	}

	fmt.Printf("Updated Entry:\nID: %d\nTitle: %s\nContent: %s\nType: %s\nCreated: %v\nUpdated: %v\n",
		journal.ID, journal.Title, journal.Content, journal.JournalType, journal.CreatedAt, journal.UpdatedAt)

	fmt.Println("\nDatabase test completed successfully!")
}
