# Agent Guidelines for test-go-htmx

## Dev process

You must follow these rules:
- Each time you make a big change to the go code, try to recompile it
- in the schema, create tables only if they don't exist.
- NEVER RUN THE APP YOURSELF

## Build/Test Commands
- **Build**: `templ generate && go build ./...` (builds all packages)
- **Generate templates**: `templ generate`

## Project Structure
- `main.go` - Application entry point
- `internal/` - Contains private application code:
  - `handlers/` - HTTP handlers for processing requests and rendering responses.
  - `models/` - Business logic, domain models, and viewmodels for passing information to templates.
  - `database/` - Low-level database connection and initialization logic.
  - `templates/` - Type-safe HTML templates for rendering views.
  - `viewmodels/` - Data structures for passing information to templates.
- `tools/` - Project-specific tools for tasks like database resets and testing.
- `web/static/` - Static assets such as CSS, JavaScript, and images.
- `data/` - SQLite database storage for application data.

## Code Style Guidelines
- **Imports**: Standard library first, then third-party, then local packages with blank lines between groups
- **Naming**: Use camelCase for variables/functions, PascalCase for exported types/functions
- **Error handling**: Always check and handle errors explicitly, use `fmt.Errorf` for wrapping
- **Logging**: Use `log.Printf` for debugging with descriptive messages
- **Database**: Use `sql.NullString` for nullable fields, always defer `rows.Close()`
- **HTTP handlers**: Check request methods explicitly, use proper HTTP status codes
- **Templates**: Use templ for type-safe HTML templates, convert models to viewmodels
- **Comments**: Document exported functions and types, explain complex logic

## Template syntax

```go
package templates

import "test-go-htmx/internal/viewmodels"

import "time"

templ ChildrenPage(children []viewmodels.Value) {
	@Base("Children | Journal App", time.Now().Year()) {
		<div>
			<h1>Children</h1>
			<table>
				<tr>
					<th>ID</th>
					<th>Name</th>
					<th>Description</th>
				</tr>
				for _, it := range children {
					<tr>
						<td>{ it.ID }</td>
						<td>{ it.Name }</td>
						<td>{ it.Description }</td>
					</tr>
				}
			</table>
		</div>
	}
}
```

## Best Practices and Code Examples

### Template Attribute Expressions
Always use expressions without quotes in HTML attributes:

```go
// Correct
<input type="hidden" name="valueID" value={ strconv.FormatInt(it.ID, 10) }/>
<option value={ strconv.FormatInt(it.ID, 10) }>{ it.Name }</option>

// Wrong
<input type="hidden" name="valueID" value="{ strconv.FormatInt(it.ID, 10) }"/>
```

### Converting int64 to String in Templates
Import `strconv` and use `FormatInt` for numeric conversions:

```go
import (
    "test-go-htmx/internal/viewmodels"
    "time"
    "strconv"
)

templ ValuesPage(values []viewmodels.Value) {
    for _, it := range values {
        <td>{ strconv.FormatInt(it.ID, 10) }</td>
    }
}
```

### HTML Form DELETE Operations
Since HTML forms only support GET and POST, create dedicated POST endpoints for deletions:

```go
// In main.go
http.HandleFunc("/values/delete", handlers.DeleteValueHandler)

// In template
<form method="POST" action="/values/delete">
    <input type="hidden" name="valueID" value={ strconv.FormatInt(it.ID, 10) }/>
    <button type="submit">Delete</button>
</form>

// In handler
func DeleteValueHandler(w http.ResponseWriter, r *http.Request) {
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

    valueIDStr := r.PostForm.Get("valueID")
    // ... rest of handler logic
}
```

### Database Operations Structure
Organize database functions properly:

```go
// internal/database/values.go - Low-level database operations
func DeleteValue(id int) error {
    query := "DELETE FROM values WHERE id = ?"
    result, err := DB.Exec(query, id)
    // ... error handling
}

// internal/models/value.go - Business logic operations
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
```


## Model Convention

The project uses a simplified model approach with a single struct per entity. SQL nullable types are converted to Go types directly in the data access functions:

```go
// Model struct
type Value struct {
    ID          int64
    Name        string
    Description string
    ParentNames string
    ParentIDs   []int64
}

// Data access function handles SQL nullable types internally
func GetAllValues() ([]Value, error) {
    // ...
    for rows.Next() {
        var v Value
        var description sql.NullString
        rows.Scan(&v.ID, &v.Name, &description)
        v.Description = description.String  // Convert SQL type to Go type
        values = append(values, v)
    }
    // ...
}
```


## Dependencies
- `github.com/mattn/go-sqlite3` - SQLite driver
