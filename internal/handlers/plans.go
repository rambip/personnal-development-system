package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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
		HandleCreatePlan(w, r)
	case http.MethodDelete:
		log.Printf("Handling DELETE request for Plans page")
		HandleDeletePlan(w, r)
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

// HandleCreatePlan creates a new plan
func HandleCreatePlan(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("name")
	description := r.PostForm.Get("description")
	resourcesRequired := r.PostForm.Get("resources") // From the form field name="resources"
	valueIDStr := r.PostForm.Get("valueID")

	// Debug logging
	log.Printf("Form data: name=%s, description=%s, resources=%s, valueID=%s",
		name, description, resourcesRequired, valueIDStr)
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

// HandleDeletePlan deletes a plan by ID
func HandleDeletePlan(w http.ResponseWriter, r *http.Request) {
	// Extract plan ID from URL or query
	var planIDStr string

	if strings.HasPrefix(r.URL.Path, "/plans/delete/") {
		// For DELETE requests (HTMX), extract from URL path
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) >= 3 {
			planIDStr = parts[len(parts)-1]
		}
	} else {
		// For traditional form submissions
		planIDStr = r.URL.Query().Get("planID")
	}

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

	// HTMX request handling - just return empty response to remove the row
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Traditional form handling
	http.Redirect(w, r, "/plans", http.StatusSeeOther)
}

// EditPlanHandler handles rendering the edit form for a plan
func EditPlanHandler(w http.ResponseWriter, r *http.Request) {
	// Extract plan ID from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	planIDStr := parts[len(parts)-1]
	planID, err := strconv.ParseInt(planIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid plan ID", http.StatusBadRequest)
		return
	}

	// Get the plan
	plan, err := models.GetPlan(planID)
	if err != nil {
		log.Printf("Error retrieving plan: %v", err)
		http.Error(w, "Error retrieving plan", http.StatusInternalServerError)
		return
	}

	// Get all values for the dropdown
	values, err := models.GetAllValues()
	if err != nil {
		log.Printf("Error retrieving values: %v", err)
		http.Error(w, "Error retrieving values", http.StatusInternalServerError)
		return
	}

	// Generate options HTML for dropdown
	valueOptions := ""
	for _, value := range values {
		selected := ""
		if value.ID == plan.ValueID {
			selected = " selected"
		}
		valueOptions += fmt.Sprintf("<option value=\"%d\"%s>%s</option>", value.ID, selected, value.Name)
	}

	// Render the edit form
	html := fmt.Sprintf(`
		<tr id="plan-row-%d" class="editing">
			<td>%d</td>
			<td>
				<input type="text" name="name" value="%s" required style="width: 100%%; box-sizing: border-box; padding: 4px;"/>
			</td>
			<td>
				<textarea name="description" required style="width: 100%%; box-sizing: border-box; padding: 4px; min-height: 60px;">%s</textarea>
			</td>
			<td>
				<input type="text" name="resources_required" value="%s" required style="width: 100%%; box-sizing: border-box; padding: 4px;"/>
			</td>
			<td>
				<select name="value_id" required style="width: 100%%; box-sizing: border-box; padding: 4px;">
					%s
				</select>
			</td>
			<td>
				<div style="display: flex; gap: 5px;">
					<button 
						hx-put="/plans/update/%d"
						hx-include="closest tr"
						hx-target="#plan-row-%d"
						hx-swap="outerHTML"
						style="background-color: #4CAF50; color: white; border: none; padding: 5px 10px; cursor: pointer; border-radius: 3px;">
						Save
					</button>
					<button 
						hx-get="/plans/cancel-edit/%d"
						hx-target="#plan-row-%d"
						hx-swap="outerHTML"
						style="background-color: #f44336; color: white; border: none; padding: 5px 10px; cursor: pointer; border-radius: 3px;">
						Cancel
					</button>
				</div>
			</td>
		</tr>
	`, plan.ID, plan.ID, plan.Name, plan.Description, plan.ResourcesRequired, valueOptions, plan.ID, plan.ID, plan.ID, plan.ID)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// CancelEditHandler handles cancelling an edit operation
func CancelEditHandler(w http.ResponseWriter, r *http.Request) {
	// Extract plan ID from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	planIDStr := parts[len(parts)-1]
	planID, err := strconv.ParseInt(planIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid plan ID", http.StatusBadRequest)
		return
	}

	// Get the plan
	plan, err := models.GetPlan(planID)
	if err != nil {
		log.Printf("Error retrieving plan: %v", err)
		http.Error(w, "Error retrieving plan", http.StatusInternalServerError)
		return
	}

	// Get all values for display
	values, err := models.GetAllValues()
	if err != nil {
		log.Printf("Error retrieving values: %v", err)
		http.Error(w, "Error retrieving values", http.StatusInternalServerError)
		return
	}

	// Find the value name
	valueName := "Unknown"
	for _, value := range values {
		if value.ID == plan.ValueID {
			valueName = value.Name
			break
		}
	}

	// Render the normal row
	html := fmt.Sprintf(`
		<tr id="plan-row-%d">
			<td>%d</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>
				<button 
					hx-get="/plans/edit/%d"
					hx-target="#plan-row-%d"
					hx-swap="outerHTML">
					Edit
				</button>
				<button
					hx-delete="/plans/delete/%d"
					hx-confirm="Are you sure you want to delete this plan?"
					hx-target="#plan-row-%d"
					hx-swap="outerHTML">
					Delete
				</button>
			</td>
		</tr>
	`, plan.ID, plan.ID, plan.Name, plan.Description, plan.ResourcesRequired, valueName, plan.ID, plan.ID, plan.ID, plan.ID)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// UpdatePlanHandler handles updating a plan
func UpdatePlanHandler(w http.ResponseWriter, r *http.Request) {
	// Extract plan ID from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	planIDStr := parts[len(parts)-1]
	planID, err := strconv.ParseInt(planIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid plan ID", http.StatusBadRequest)
		return
	}

	// Parse form
	err = r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get form values
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

	// Log the form values for debugging
	log.Printf("Updating plan - ID: %d, Name: %s, Description: %s, Resources Required: %s, Value ID: %d",
		planID, name, description, resourcesRequired, valueID)

	// Update the plan
	err = models.UpdatePlan(planID, name, description, resourcesRequired, valueID)
	if err != nil {
		log.Printf("Error updating plan: %v", err)
		http.Error(w, "Error updating plan", http.StatusInternalServerError)
		return
	}

	// Get the updated plan
	plan, err := models.GetPlan(planID)
	if err != nil {
		log.Printf("Error retrieving updated plan: %v", err)
		http.Error(w, "Error retrieving updated plan", http.StatusInternalServerError)
		return
	}

	// Get all values for display
	values, err := models.GetAllValues()
	if err != nil {
		log.Printf("Error retrieving values: %v", err)
		http.Error(w, "Error retrieving values", http.StatusInternalServerError)
		return
	}

	// Find the value name
	valueName := "Unknown"
	for _, value := range values {
		if value.ID == plan.ValueID {
			valueName = value.Name
			break
		}
	}

	// Render the updated row
	html := fmt.Sprintf(`
		<tr id="plan-row-%d">
			<td>%d</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>
				<button 
					hx-get="/plans/edit/%d"
					hx-target="#plan-row-%d"
					hx-swap="outerHTML">
					Edit
				</button>
				<button
					hx-delete="/plans/delete/%d"
					hx-confirm="Are you sure you want to delete this plan?"
					hx-target="#plan-row-%d"
					hx-swap="outerHTML">
					Delete
				</button>
			</td>
		</tr>
	`, plan.ID, plan.ID, plan.Name, plan.Description, plan.ResourcesRequired, valueName, plan.ID, plan.ID, plan.ID, plan.ID)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
