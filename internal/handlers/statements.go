package handlers

import (
	"log"
	"net/http"
	"strconv"
	"test-go-htmx/internal/models"
	"test-go-htmx/internal/templates"
)

// StatementsHandler handles the Statements page
func StatementsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("StatementsHandler called with path: %s, method: %s", r.URL.Path, r.Method)

	// Handle GET requests
	if r.Method == http.MethodGet {
		log.Printf("Handling GET request for Statements page")
		handleGetStatements(w, r)
	} else {
		log.Printf("Method %s not allowed for Statements page", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// CreateStatementHandler handles creating new statements
func CreateStatementHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	content := r.PostForm.Get("content")
	priorityStr := r.PostForm.Get("priority")
	priority, err := strconv.Atoi(priorityStr)
	if err != nil {
		log.Printf("Invalid priority: %v", err)
		priority = 0 // Default priority
	}

	log.Printf("Creating new statement - Content: %s, Priority: %d", content, priority)

	id, err := models.CreateStatement(content, priority)
	if err != nil {
		log.Printf("Error creating statement: %v", err)
		http.Error(w, "Error creating statement", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully created statement with ID: %d", id)
	http.Redirect(w, r, "/statements", http.StatusSeeOther)
}

// DeleteStatementHandler handles deleting statements
func DeleteStatementHandler(w http.ResponseWriter, r *http.Request) {
	var statementIDStr string

	if r.Method == http.MethodPost {
		// Parse form for POST requests
		if err := r.ParseForm(); err != nil {
			log.Printf("Error parsing form: %v", err)
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		statementIDStr = r.PostForm.Get("statementID")
		log.Printf("POST Delete - Form data: %+v, statementID: %s", r.PostForm, statementIDStr)
	} else if r.Method == http.MethodGet {
		// Get query parameters for GET requests
		statementIDStr = r.URL.Query().Get("statementID")
		log.Printf("GET Delete - Query params: %+v, statementID: %s", r.URL.Query(), statementIDStr)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the statement ID
	statementID, err := strconv.ParseInt(statementIDStr, 10, 64)
	if err != nil {
		log.Printf("Invalid statement ID: %v", err)
		http.Error(w, "Invalid statement ID", http.StatusBadRequest)
		return
	}

	if statementID <= 0 {
		log.Printf("Missing statement ID")
		http.Error(w, "statementID is required", http.StatusBadRequest)
		return
	}

	// Delete the statement
	err = models.DeleteStatement(statementID)
	if err != nil {
		log.Printf("Error deleting statement: %v", err)
		http.Error(w, "Error deleting statement", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully deleted statement with ID: %d", statementID)
	http.Redirect(w, r, "/statements", http.StatusSeeOther)
}

// handleGetStatements retrieves and displays all statements
func handleGetStatements(w http.ResponseWriter, r *http.Request) {
	statements, err := models.GetAllStatements()
	if err != nil {
		log.Printf("Error retrieving statements: %v", err)
		http.Error(w, "Error retrieving statements", http.StatusInternalServerError)
		return
	}

	component := templates.StatementsPage(statements)
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering Statements page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully rendered Statements page with %d statements", len(statements))
}
