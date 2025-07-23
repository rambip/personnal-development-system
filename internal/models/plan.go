package models

import (
	"pds/internal/database"
)

// Plan represents a plan in the system
type Plan struct {
	ID                int64
	Name              string
	Description       string
	ResourcesRequired string
	ValueID           int64
}

// CreatePlan inserts a new plan into the database
func CreatePlan(name, description, resourcesRequired string, valueID int64) (int64, error) {
	query := "INSERT INTO plans (name, description, resources_required, value_id) VALUES (?, ?, ?, ?)"
	result, err := database.DB.Exec(query, name, description, resourcesRequired, valueID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetPlan retrieves a plan by ID
func GetPlan(id int64) (Plan, error) {
	query := "SELECT id, name, description, resources_required, value_id FROM plans WHERE id = ?"
	var plan Plan
	row := database.DB.QueryRow(query, id)
	if err := row.Scan(&plan.ID, &plan.Name, &plan.Description, &plan.ResourcesRequired, &plan.ValueID); err != nil {
		return plan, err
	}
	return plan, nil
}

// GetAllPlans retrieves all plans from the database
func GetAllPlans() ([]Plan, error) {
	query := "SELECT id, name, description, resources_required, value_id FROM plans"
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []Plan
	for rows.Next() {
		var plan Plan
		if err := rows.Scan(&plan.ID, &plan.Name, &plan.Description, &plan.ResourcesRequired, &plan.ValueID); err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	return plans, nil
}

// UpdatePlan updates an existing plan
func UpdatePlan(id int64, name, description, resourcesRequired string, valueID int64) error {
	query := "UPDATE plans SET name = ?, description = ?, resources_required = ?, value_id = ? WHERE id = ?"
	_, err := database.DB.Exec(query, name, description, resourcesRequired, valueID, id)
	return err
}

// DeletePlan deletes a plan by ID
func DeletePlan(id int64) error {
	query := "DELETE FROM plans WHERE id = ?"
	_, err := database.DB.Exec(query, id)
	return err
}
