package handlers

import (
	"log"
	"net/http"
	"strconv"
	"test-go-htmx/internal/models"
	"test-go-htmx/internal/templates"
)

// PlansHandler handles the Plans page
func PlansHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("PlansHandler called with path: %s, method: %s", r.URL.Path, r.Method)

	// Handle GET, POST, and DELETE requests
	switch r.Method {
	case http.MethodGet:
		log.Printf("Handling GET request for Plans page")
		handleGetPlans(w, r)
	case http.MethodPost:
		log.Printf("Handling POST request for Plans page")
		handleCreatePlan(w, r)
	case http.MethodDelete:
		log.Printf("Handling DELETE request for Plans page")
		handleDeletePlan(w, r)
	default:
		log.Printf("Method %s not allowed for Plans page", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetPlans retrieves and displays all plans
func handleGetPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := models.GetAllPlans()
	values, err := models.GetAllValues()
	if err != nil {
		log.Printf("Error retrieving values: %v", err)
		http.Error(w, "Error retrieving values", http.StatusInternalServerError)
		return
	}
	if err != nil {
		log.Printf("Error retrieving plans: %v", err)
		http.Error(w, "Error retrieving plans", http.StatusInternalServerError)
		return
	}

	component := templates.PlansPage(plans, values)
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering Plans page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully rendered Plans page")
}

// handleCreatePlan creates a new plan
func handleCreatePlan(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("name")
	description := r.PostForm.Get("description")
	resourcesRequired := r.PostForm.Get("resources_required")
	valueIDStr := r.PostForm.Get("value_id")
	valueID, err := strconv.ParseInt(valueIDStr, 10, 64)
	if err != nil {
		log.Printf("Invalid value ID: %v", err)
		http.Error(w, "Invalid value ID", http.StatusBadRequest)
		return
	}

	log.Printf("Creating new plan - Name: %s, Description: %s, Resources Required: %s, Value ID: %d", name, description, resourcesRequired, valueID)

	id, err := models.CreatePlan(name, description, resourcesRequired, valueID)
	if err != nil {
		log.Printf("Error creating plan: %v", err)
		http.Error(w, "Error creating plan", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully created plan with ID: %d", id)
	http.Redirect(w, r, "/plans", http.StatusSeeOther)
}

// handleDeletePlan deletes a plan by ID
func handleDeletePlan(w http.ResponseWriter, r *http.Request) {
	planIDStr := r.URL.Query().Get("planID")
	planID, err := strconv.ParseInt(planIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid plan ID", http.StatusBadRequest)
		return
	}
	if planID <= 0 {
		http.Error(w, "planID is required", http.StatusBadRequest)
		return
	}

	err = models.DeletePlan(planID)
	if err != nil {
		log.Printf("Error deleting plan: %v", err)
		http.Error(w, "Error deleting plan", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully deleted plan with ID: %d", planID)
	http.Redirect(w, r, "/plans", http.StatusSeeOther)
}
