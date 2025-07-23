package handlers

import (
	"log"
	"net/http"
	"strconv"
	"test-go-htmx/internal/models"
	"test-go-htmx/internal/templates"
	"test-go-htmx/internal/viewmodels"
)

// ValuesHandler handles the Values page.

// handleGetChildren retrieves and displays children for a specific value.
func handleGetChildren(w http.ResponseWriter, r *http.Request) {
	valueIDStr := r.URL.Query().Get("valueID")
	valueID, err := strconv.ParseInt(valueIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid valueID", http.StatusBadRequest)
		return
	}
	if valueID <= 0 {
		http.Error(w, "valueID is required", http.StatusBadRequest)
		return
	}

	children, err := models.GetChildren(valueID)
	if err != nil {
		log.Printf("Error retrieving children: %v", err)
		http.Error(w, "Error retrieving children", http.StatusInternalServerError)
		return
	}

	viewChildren := viewmodels.ConvertModelsToViewValues(children)
	component := templates.ChildrenPage(viewChildren)
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering Children page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully rendered Children page")
}

// handleGetParents retrieves and displays parents for a specific value.
func handleGetParents(w http.ResponseWriter, r *http.Request) {
	valueIDStr := r.URL.Query().Get("valueID")
	valueID, err := strconv.ParseInt(valueIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid valueID", http.StatusBadRequest)
		return
	}
	if valueID <= 0 {
		http.Error(w, "valueID is required", http.StatusBadRequest)
		return
	}

	parents, err := models.GetParents(valueID)
	if err != nil {
		log.Printf("Error retrieving parents: %v", err)
		http.Error(w, "Error retrieving parents", http.StatusInternalServerError)
		return
	}

	viewParents := viewmodels.ConvertModelsToViewValues(parents)
	component := templates.ParentsPage(viewParents)
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering Parents page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully rendered Parents page")
}

func ValuesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("ValuesHandler called with path: %s, method: %s", r.URL.Path, r.Method)

	// Handle GET, POST, and DELETE requests
	switch r.Method {
	case http.MethodGet:
		log.Printf("Handling GET request for Values page")
		handleGetValues(w, r)
	case http.MethodPost:
		log.Printf("Handling POST request for Values page")
		handleCreateValue(w, r)
	case http.MethodDelete:
		log.Printf("Handling DELETE request for Values page")
		handleDeleteValue(w, r)
	default:
		log.Printf("Method %s not allowed for Values page", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetValues retrieves and displays all values.
func handleGetValues(w http.ResponseWriter, r *http.Request) {
	values, err := models.GetAllValues()
	if err != nil {
		log.Printf("Error retrieving values: %v", err)
		http.Error(w, "Error retrieving values", http.StatusInternalServerError)
		return
	}

	viewValues := viewmodels.ConvertModelsToViewValues(values)
	component := templates.ValuesPage(viewValues)
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering Values page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully rendered Values page")
}

// handleCreateValue creates a new value with parent relationships.
func handleCreateValue(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("name")
	description := r.PostForm.Get("description")
	parentIDs := r.PostForm["parents"]

	log.Printf("Creating new value - Name: %s, Description: %s, Parent IDs: %v", name, description, parentIDs)

	if name == "" {
		log.Printf("Validation failed: name is empty")
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	id, err := models.CreateValue(name, description, parentIDs)
	if err != nil {
		log.Printf("Error creating value: %v", err)
		http.Error(w, "Error creating value", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully created value with ID: %d", id)
	http.Redirect(w, r, "/values", http.StatusSeeOther)
}

// handleDeleteValue deletes a value by ID.
func handleDeleteValue(w http.ResponseWriter, r *http.Request) {
	valueIDStr := r.URL.Query().Get("valueID")
	valueID, err := strconv.ParseInt(valueIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid valueID", http.StatusBadRequest)
		return
	}
	if valueID <= 0 {
		http.Error(w, "valueID is required", http.StatusBadRequest)
		return
	}

	err = models.DeleteValue(valueID)
	if err != nil {
		log.Printf("Error deleting value: %v", err)
		http.Error(w, "Error deleting value", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully deleted value with ID: %d", valueID)
	http.Redirect(w, r, "/values", http.StatusSeeOther)
}
