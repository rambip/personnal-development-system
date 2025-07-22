package viewmodels

import (
	"time"
)

// JournalEntry represents a journal entry for template rendering
type JournalEntry struct {
	ID          int64
	Title       string
	Content     string
	JournalType string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// FormatTime formats a time for display in templates
func FormatTime(t time.Time) string {
	return t.Format("Jan 02, 2006 at 15:04")
}
