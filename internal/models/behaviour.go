package models

import (
	"pds/internal/database"
)

// Behaviour represents a behaviour that conflicts with an aim
type Behaviour struct {
	ID                 int64
	Name               string
	Description        string
	Mark               string
	ConflictingAimID   int64
	ConflictingAimName string // For display purposes
}

// CreateBehaviour inserts a new behaviour into the database
func CreateBehaviour(name, description, mark string, conflictingAimID int64) (int64, error) {
	query := "INSERT INTO behaviours (name, description, mark, conflicting_aim_id) VALUES (?, ?, ?, ?)"
	result, err := database.DB.Exec(query, name, description, mark, conflictingAimID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetAllBehaviours retrieves all behaviours with their conflicting aim names
func GetAllBehaviours() ([]Behaviour, error) {
	query := `
		SELECT b.id, b.name, b.description, b.mark, b.conflicting_aim_id, a.name
		FROM behaviours b
		LEFT JOIN aims a ON b.conflicting_aim_id = a.id
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var behaviours []Behaviour
	for rows.Next() {
		var behaviour Behaviour
		if err := rows.Scan(
			&behaviour.ID,
			&behaviour.Name,
			&behaviour.Description,
			&behaviour.Mark,
			&behaviour.ConflictingAimID,
			&behaviour.ConflictingAimName,
		); err != nil {
			return nil, err
		}
		behaviours = append(behaviours, behaviour)
	}
	return behaviours, nil
}

// DeleteBehaviour deletes a behaviour by ID
func DeleteBehaviour(id int64) error {
	query := "DELETE FROM behaviours WHERE id = ?"
	_, err := database.DB.Exec(query, id)
	return err
}
