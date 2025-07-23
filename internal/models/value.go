package models

import (
	"database/sql"
	"test-go-htmx/internal/database"
)

// Ensure Value type is defined
type Value struct {
	ID          int64
	Name        string
	Description sql.NullString
	ParentIDs   []int64
}

// GetChildren retrieves all child values for a given value ID.
func GetChildren(valueID int64) ([]Value, error) {
	db := database.DB
	rows, err := db.Query(
		`SELECT v.id, v.name, v.description
		 FROM values v
		 JOIN value_parents vp ON v.id = vp.value_id
		 WHERE vp.parent_value_id = ?`,
		valueID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []Value
	for rows.Next() {
		var v Value
		var description sql.NullString
		if err := rows.Scan(&v.ID, &v.Name, &description); err != nil {
			return nil, err
		}
		v.Description = description
		children = append(children, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return children, nil
}

// GetParents retrieves all parent values for a given value ID.
func GetParents(valueID int64) ([]Value, error) {
	db := database.DB
	rows, err := db.Query(
		`SELECT v.id, v.name, v.description
		 FROM values v
		 JOIN value_parents vp ON v.id = vp.parent_value_id
		 WHERE vp.value_id = ?`,
		valueID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parents []Value
	for rows.Next() {
		var v Value
		var description sql.NullString
		if err := rows.Scan(&v.ID, &v.Name, &description); err != nil {
			return nil, err
		}
		v.Description = description
		parents = append(parents, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return parents, nil
}

// GetAllValues retrieves all values from the database.
func GetAllValues() ([]Value, error) {
	db := database.DB
	rows, err := db.Query("SELECT id, name, description FROM `values`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []Value
	for rows.Next() {
		var v Value
		var description sql.NullString
		if err := rows.Scan(&v.ID, &v.Name, &description); err != nil {
			return nil, err
		}
		v.Description = description
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
		"INSERT INTO `values` (name, description) VALUES (?, ?)",
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
		return err
	}

	// Delete the value itself
	_, err = db.Exec("DELETE FROM `values` WHERE id = ?", valueID)
	return err
}
