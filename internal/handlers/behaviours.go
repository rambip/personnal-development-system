package handlers

import (
	"log"
	"net/http"
	"strconv"
	"test-go-htmx/internal/models"
	"test-go-htmx/internal/templates"
)

// BehavioursHandler handles the Behaviours page
func BehavioursHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("BehavioursHandler called with path: %s, method: %s", r.URL.Path, r.Method)

	// Handle GET requests
	if r.Method == http.MethodGet {
		log.Printf("Handling GET request for Behaviours page")
		handleGetBehaviours(w, r)
	} else {
		log.Printf("Method %s not allowed for Behaviours page", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// CreateBehaviourHandler handles creating new behaviours
func CreateBehaviourHandler(w http.ResponseWriter, r *http.Request) {
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

	name := r.PostForm.Get("name")
	description := r.PostForm.Get("description")
	mark := r.PostForm.Get("mark")
	conflictingAimIDStr := r.PostForm.Get("conflictingAimID")

	conflictingAimID, err := strconv.ParseInt(conflictingAimIDStr, 10, 64)
	if err != nil {
		log.Printf("Invalid conflicting aim ID: %v", err)
		http.Error(w, "Invalid conflicting aim ID", http.StatusBadRequest)
		return
	}

	log.Printf("Creating new behaviour - Name: %s, Description: %s, Mark: %s, Conflicting Aim ID: %d",
		name, description, mark, conflictingAimID)

	id, err := models.CreateBehaviour(name, description, mark, conflictingAimID)
	if err != nil {
		log.Printf("Error creating behaviour: %v", err)
		http.Error(w, "Error creating behaviour", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully created behaviour with ID: %d", id)

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// Return just the updated behaviours list for HTMX
		behaviours, err := models.GetAllBehaviours()
		if err != nil {
			log.Printf("Error retrieving behaviours: %v", err)
			http.Error(w, "Error retrieving behaviours", http.StatusInternalServerError)
			return
		}

		component := templates.BehavioursList(behaviours)
		if err := component.Render(r.Context(), w); err != nil {
			log.Printf("Error rendering behaviours list: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	} else {
		// Regular form submit - redirect to behaviours page
		http.Redirect(w, r, "/behaviours", http.StatusSeeOther)
	}
}

// DeleteBehaviourHandler handles deleting behaviours
func DeleteBehaviourHandler(w http.ResponseWriter, r *http.Request) {
	// Support both POST and DELETE methods (HTMX uses DELETE)
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var behaviourIDStr string

	// Parse the appropriate data based on the request method
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			log.Printf("Error parsing form: %v", err)
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		behaviourIDStr = r.PostForm.Get("behaviourID")
		log.Printf("POST Delete - Form data: %+v", r.PostForm)
	} else { // DELETE method
		if err := r.ParseForm(); err != nil {
			log.Printf("Error parsing form: %v", err)
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		behaviourIDStr = r.Form.Get("behaviourID")
		log.Printf("DELETE request - Form data: %+v", r.Form)
	}

	log.Printf("Deleting behaviour with ID: %s", behaviourIDStr)

	// Parse the behaviour ID
	behaviourID, err := strconv.ParseInt(behaviourIDStr, 10, 64)
	if err != nil {
		log.Printf("Invalid behaviour ID: %v", err)
		http.Error(w, "Invalid behaviour ID", http.StatusBadRequest)
		return
	}

	if behaviourID <= 0 {
		log.Printf("Missing behaviour ID")
		http.Error(w, "behaviourID is required", http.StatusBadRequest)
		return
	}

	// Delete the behaviour
	err = models.DeleteBehaviour(behaviourID)
	if err != nil {
		log.Printf("Error deleting behaviour: %v", err)
		http.Error(w, "Error deleting behaviour", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully deleted behaviour with ID: %d", behaviourID)

	// Handle HTMX request differently
	if r.Header.Get("HX-Request") == "true" {
		// For HTMX DELETE requests, we can just return an empty response with 200 status
		// This will remove the row from the table due to hx-swap="outerHTML"
		w.WriteHeader(http.StatusOK)
	} else {
		// Regular form submit - redirect to behaviours page
		http.Redirect(w, r, "/behaviours", http.StatusSeeOther)
	}
}

// handleGetBehaviours retrieves and displays all behaviours
func handleGetBehaviours(w http.ResponseWriter, r *http.Request) {
	behaviours, err := models.GetAllBehaviours()
	if err != nil {
		log.Printf("Error retrieving behaviours: %v", err)
		http.Error(w, "Error retrieving behaviours", http.StatusInternalServerError)
		return
	}

	aims, err := models.GetAllValues()
	if err != nil {
		log.Printf("Error retrieving aims: %v", err)
		http.Error(w, "Error retrieving aims", http.StatusInternalServerError)
		return
	}

	component := templates.BehavioursPage(behaviours, aims)
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering Behaviours page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully rendered Behaviours page with %d behaviours", len(behaviours))
}
