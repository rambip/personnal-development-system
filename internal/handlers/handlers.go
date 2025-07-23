package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"pds/internal/models"
	"pds/internal/templates"
)

// These conversion functions are no longer needed with the simplified model approach

// HomeHandler handles the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("HomeHandler called with path: %s", r.URL.Path)
	if r.URL.Path != "/" {
		log.Printf("Path %s is not '/', returning 404", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	component := templates.Home()
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering home template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully rendered home template")
}

// JournalsHandler handles the journals page
func JournalsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("JournalsHandler called with path: %s, method: %s", r.URL.Path, r.Method)

	// Different behavior based on HTTP method
	switch r.Method {
	case http.MethodGet:
		log.Printf("Handling GET request for journals")
		handleGetJournals(w, r)
	case http.MethodPost:
		log.Printf("Handling POST request for journals")
		handleCreateJournal(w, r)
	default:
		log.Printf("Method %s not allowed for journals", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetJournals handles GET requests for journal entries
func handleGetJournals(w http.ResponseWriter, r *http.Request) {
	log.Printf("handleGetJournals with path: %s", r.URL.Path)
	var journals []models.Journal
	var err error

	// Check if we're filtering by type
	if strings.Contains(r.URL.Path, "/type/") {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) >= 4 {
			journalType := parts[3]
			log.Printf("Filtering journals by type: %s", journalType)
			journals, err = models.GetJournalsByType(journalType)
		}
	} else {
		log.Printf("Retrieving all journals")
		journals, err = models.GetAllJournals()
	}

	if err != nil {
		log.Printf("Error retrieving journals: %v", err)
		http.Error(w, "Error retrieving journals", http.StatusInternalServerError)
		return
	}

	log.Printf("Retrieved %d journals", len(journals))

	// If it's an HTMX request, just return the journal list partial
	if r.Header.Get("HX-Request") == "true" {
		log.Printf("HTMX request detected, rendering partial template")
		component := templates.JournalList(journals)
		if err := component.Render(r.Context(), w); err != nil {
			log.Printf("Error rendering partial template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Printf("Successfully rendered partial template")
		return
	}

	// Otherwise, return the full page
	log.Printf("Rendering full journals page")
	component := templates.Journals(journals)
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering journals template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully rendered journals template")
}

// handleCreateJournal handles POST requests to create a new journal entry
func handleCreateJournal(w http.ResponseWriter, r *http.Request) {
	log.Printf("handleCreateJournal called")
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Extract form values
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	journalType := r.PostForm.Get("journal_type")
	log.Printf("Creating new journal entry - Title: %s, Type: %s, Content length: %d",
		title, journalType, len(content))

	// Validate form values
	if title == "" || journalType == "" {
		log.Printf("Validation failed: title or journal_type is empty")
		http.Error(w, "Title and journal type are required", http.StatusBadRequest)
		return
	}

	// Create journal entry
	id, err := models.CreateJournal(title, content, journalType)
	if err != nil {
		log.Printf("Error creating journal: %v", err)
		http.Error(w, "Error creating journal", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully created journal with ID: %d", id)

	// Get the newly created journal entry
	journal, err := models.GetJournal(id)
	if err != nil {
		log.Printf("Error retrieving created journal: %v", err)
		http.Error(w, "Error retrieving created journal", http.StatusInternalServerError)
		return
	}
	log.Printf("Retrieved created journal: ID=%d, Title=%s, Type=%s",
		journal.ID, journal.Title, journal.JournalType)

	// Return just the single journal entry if it's an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		log.Printf("Responding to HTMX create request with partial template")
		component := templates.JournalEntry(journal)
		if err := component.Render(r.Context(), w); err != nil {
			log.Printf("Error rendering partial template after create: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Printf("Successfully rendered partial template after create")
		return
	}

	// Redirect to journals page if it's not an HTMX request
	log.Printf("Redirecting to journals page after create")
	http.Redirect(w, r, "/journals", http.StatusSeeOther)
}

// HandleDeleteJournal handles POST requests to delete a journal entry
func HandleDeleteJournal(w http.ResponseWriter, r *http.Request) {
	log.Printf("handleDeleteJournal called")
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Extract journal ID from form
	idStr := r.PostForm.Get("journalID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Invalid journal ID: %s - %v", idStr, err)
		http.Error(w, "Invalid journal ID", http.StatusBadRequest)
		return
	}

	// Delete the journal entry
	err = models.DeleteJournal(id)
	if err != nil {
		log.Printf("Error deleting journal: %v", err)
		http.Error(w, "Error deleting journal", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully deleted journal with ID: %d", id)

	// Respond to HTMX request
	if r.Header.Get("HX-Request") == "true" {
		log.Printf("Responding to HTMX delete request")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Redirect to journals page if not HTMX
	http.Redirect(w, r, "/journals", http.StatusSeeOther)
}

// JournalDetailHandler handles requests for a specific journal entry
func JournalDetailHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("JournalDetailHandler called with path: %s", r.URL.Path)

	// Extract journal ID from URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		log.Printf("Invalid path format for journal detail: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		log.Printf("Invalid journal ID: %s - %v", parts[2], err)
		http.NotFound(w, r)
		return
	}
	log.Printf("Requested journal with ID: %d", id)

	// Just check if the journal exists
	_, err = models.GetJournal(id)
	if err != nil {
		log.Printf("Error retrieving journal: %v", err)
		http.Error(w, "Error retrieving journal", http.StatusInternalServerError)
		return
	}

	// This handler could be expanded to show a detailed view of a journal entry
	// For now, we'll just redirect to the journals page
	http.Redirect(w, r, "/journals", http.StatusSeeOther)
}
