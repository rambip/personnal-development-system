# Agent Guidelines for test-go-htmx

## Dev process

You must follow these rules:
- Each time you make a big change to the go code, try to recompile it
- in the schema, create tables only if they don't exist.
- NEVER RUN THE APP YOURSELF

## Build/Test Commands
- **Build**: `go build ./...` (builds all packages)
- **Generate templates**: `templ generate`
- **Run**: `go run main.go` (starts server on :8888)
- **Test**: `go test ./...` (runs all tests)
- **Single test**: `go test -run TestName ./path/to/package`
- **Format**: `go fmt ./...`
- **Vet**: `go vet ./...`

## Project Structure
- `main.go` - Application entry point
- `internal/` - Private application code (handlers, models, database, templates, viewmodels)
- `tools` - project specific tools
- `web/static/` - Static assets (CSS, JS, images)
- `data/` - SQLite database storage

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

## Dependencies
- `github.com/mattn/go-sqlite3` - SQLite driver
