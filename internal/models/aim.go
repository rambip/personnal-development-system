package models

import (
	"database/sql"
	"fmt"
	"pds/internal/database"
)

// Aim represents a value in the system
type Aim struct {
	ID          int64
	Name        string
	Description string
	ParentNames string
	ParentIDs   []int64
}

// GetChildren retrieves all child values for a given value ID.
func GetChildren(valueID int64) ([]Aim, error) {
	db := database.DB
	rows, err := db.Query(
		`SELECT v.id, v.name, v.description
		 FROM aims v
		 JOIN value_parents vp ON v.id = vp.value_id
		 WHERE vp.parent_value_id = ?`,
		valueID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []Aim
	for rows.Next() {
		var v Aim
		var description sql.NullString
		if err := rows.Scan(&v.ID, &v.Name, &description); err != nil {
			return nil, err
		}
		v.Description = description.String
		children = append(children, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return children, nil
}

// GetParents retrieves all parent values for a given Aim ID.
func GetParents(valueID int64) ([]Aim, error) {
	db := database.DB
	rows, err := db.Query(
		`SELECT v.id, v.name, v.description
		 FROM aims v
		 JOIN value_parents vp ON v.id = vp.parent_value_id
		 WHERE vp.value_id = ?`,
		valueID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parents []Aim
	for rows.Next() {
		var v Aim
		var description sql.NullString
		if err := rows.Scan(&v.ID, &v.Name, &description); err != nil {
			return nil, err
		}
		v.Description = description.String
		parents = append(parents, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return parents, nil
}

// GetAllValues retrieves all Aim from the database.
func GetAllValues() ([]Aim, error) {
	db := database.DB
	rows, err := db.Query("SELECT id, name, description FROM aims")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []Aim
	for rows.Next() {
		var v Aim
		var description sql.NullString
		if err := rows.Scan(&v.ID, &v.Name, &description); err != nil {
			return nil, err
		}
		v.Description = description.String
		values = append(values, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return values, nil
}

// CreateValue inserts a new value into the database.
func CreateValue(name string, description string, parentIDs []string) (int64, error) {
	db := database.DB
	result, err := db.Exec(
		"INSERT INTO aims (name, description) VALUES (?, ?)",
		name, description,
	)
	if err != nil {
		return 0, err
	}

	valueID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, parentID := range parentIDs {
		_, err := db.Exec(
			"INSERT INTO value_parents (value_id, parent_value_id) VALUES (?, ?)",
			valueID, parentID,
		)
		if err != nil {
			return 0, err
		}
	}

	return valueID, nil
}

// DeleteValue removes a value and its relationships from the database.
func DeleteValue(valueID int64) error {
	db := database.DB

	// Delete relationships first
	_, err := db.Exec("DELETE FROM value_parents WHERE value_id = ? OR parent_value_id = ?", valueID, valueID)
	if err != nil {
		return fmt.Errorf("failed to delete value relationships: %w", err)
	}

	// Delete the value itself
	result, err := db.Exec("DELETE FROM aims WHERE id = ?", valueID)
	if err != nil {
		return fmt.Errorf("failed to delete value: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no value found with ID %d", valueID)
	}

	return nil
}
