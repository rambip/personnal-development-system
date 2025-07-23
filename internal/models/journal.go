package models

import (
	"database/sql"
	"fmt"
	"time"

	"pds/internal/database"
)

// Journal represents a journal entry in the database
type Journal struct {
	ID          int64
	Title       string
	Content     string
	JournalType string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GetAllJournals retrieves all journal entries from the database
func GetAllJournals() ([]Journal, error) {
	db := database.DB
	rows, err := db.Query("SELECT id, title, content, journal_type, created_at, updated_at FROM journals ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var journals []Journal
	for rows.Next() {
		var j Journal
		var content sql.NullString
		err := rows.Scan(&j.ID, &j.Title, &content, &j.JournalType, &j.CreatedAt, &j.UpdatedAt)
		if err != nil {
			return nil, err
		}
		j.Content = content.String
		journals = append(journals, j)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return journals, nil
}

// GetJournal retrieves a journal entry by ID
func GetJournal(id int64) (Journal, error) {
	db := database.DB
	var j Journal
	var content sql.NullString
	err := db.QueryRow("SELECT id, title, content, journal_type, created_at, updated_at FROM journals WHERE id = ?", id).
		Scan(&j.ID, &j.Title, &content, &j.JournalType, &j.CreatedAt, &j.UpdatedAt)
	if err == nil {
		j.Content = content.String
	}
	return j, err
}

// CreateJournal inserts a new journal entry into the database
func CreateJournal(title string, content string, journalType string) (int64, error) {
	db := database.DB
	result, err := db.Exec(
		"INSERT INTO journals (title, content, journal_type) VALUES (?, ?, ?)",
		title, content, journalType,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// UpdateJournal updates an existing journal entry
func UpdateJournal(id int64, title string, content string, journalType string) error {
	db := database.DB
	_, err := db.Exec(
		"UPDATE journals SET title = ?, content = ?, journal_type = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		title, content, journalType, id,
	)
	return err
}

// GetJournalsByType retrieves all journal entries of a specific type
func GetJournalsByType(journalType string) ([]Journal, error) {
	db := database.DB
	rows, err := db.Query("SELECT id, title, content, journal_type, created_at, updated_at FROM journals WHERE journal_type = ? ORDER BY created_at DESC", journalType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var journals []Journal
	for rows.Next() {
		var j Journal
		var content sql.NullString
		err := rows.Scan(&j.ID, &j.Title, &content, &j.JournalType, &j.CreatedAt, &j.UpdatedAt)
		if err != nil {
			return nil, err
		}
		j.Content = content.String
		journals = append(journals, j)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return journals, nil
}

// DeleteJournal deletes a journal entry by ID
func DeleteJournal(id int64) error {
	db := database.DB
	query := "DELETE FROM journals WHERE id = ?"
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete journal entry: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no journal entry found with ID %d", id)
	}

	return nil
}
