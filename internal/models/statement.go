package models

import (
	"test-go-htmx/internal/database"
)

// Statement represents a statement in the system
type Statement struct {
	ID       int64
	Content  string
	Priority int
}

// CreateStatement inserts a new statement into the database
func CreateStatement(content string, priority int) (int64, error) {
	query := "INSERT INTO statements (content, priority) VALUES (?, ?)"
	result, err := database.DB.Exec(query, content, priority)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetAllStatements retrieves all statements from the database
func GetAllStatements() ([]Statement, error) {
	query := "SELECT id, content, priority FROM statements"
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statements []Statement
	for rows.Next() {
		var statement Statement
		if err := rows.Scan(&statement.ID, &statement.Content, &statement.Priority); err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}
	return statements, nil
}

// DeleteStatement deletes a statement by ID
func DeleteStatement(id int64) error {
	query := "DELETE FROM statements WHERE id = ?"
	_, err := database.DB.Exec(query, id)
	return err
}
