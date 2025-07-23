# A Personnal Development System

<img width="937" height="871" alt="image" src="https://github.com/user-attachments/assets/a6fb70a6-6a69-447d-996c-75e492f8d4e6" />


A Go application for tracking journals, values, plans, statements, and behaviors with HTMX-powered UI.

This is meant to be an implementation of a [Personnal Development System](https://www.lesswrong.com/posts/mpbtk2xBjqjL7p5uQ/personal-development-system-winning-repeatedly-and-growing)

## Features
- [x] Journal entries with different types
- [x] Value tracking and hierarchies
- [x] Plan management
- [x] Statement (mantras) expression
- [x] Behavior tracking
- [ ] Moods
- [ ] LLM conversation

## Project Structure
```
/test-go-htmx
├── main.go             # Application entry point
├── internal/           # Contains private application code
│   ├── database/       # Database connection and migrations
│   │   └── migrations/ # SQL migration files
│   ├── handlers/       # HTTP handlers for processing requests
│   ├── models/         # Business logic and domain models
│   └── templates/      # Type-safe templ HTML templates
├── web/                # Web assets
│   └── static/         # Static files (CSS, JS, images)
├── tools/              # Helper tools and utilities
├── go.mod              # Go module file
└── README.md           # Project documentation
```

## Quick Start Options

### Option 1: Local Development Setup

#### Prerequisites
- Go 1.18+
- [Templ](https://github.com/a-h/templ) (for generating templates)
- SQLite (included in Go's standard library via mattn/go-sqlite3)

#### Steps
1. **Clone the repository**
```bash
git clone https://github.com/yourusername/test-go-htmx.git
cd test-go-htmx
```

2. **Install Templ**
```bash
go install github.com/a-h/templ/cmd/templ@latest
```

3. **Generate templates and build**
```bash
templ generate
go build ./...
```

4. **Run the application**
```bash
./test-go-htmx
```

5. **Access the application**
Open your browser and navigate to http://localhost:8888


## License
This project is licensed under the MIT License. See the LICENSE file for details.
