package viewmodels

import "test-go-htmx/internal/models"

// Value represents a viewmodel for the Values page.
type Value struct {
	ID          int64
	Name        string
	Description string
	ParentNames string
}

// ConvertModelsToViewValues converts a slice of models.Value to viewmodels.Value.
func ConvertModelsToViewValues(values []models.Value) []Value {
	result := make([]Value, len(values))
	for i, v := range values {
		result[i] = Value{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description.String,
			ParentNames: "", // Placeholder for parent names, to be implemented
		}
	}
	return result
}
