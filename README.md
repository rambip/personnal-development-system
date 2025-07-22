# Go + HTML Project Template

## Overview
This project is a template for building web applications using Go for the backend and HTML for the frontend. It includes a basic structure for organizing code, templates, and static assets.

## Project Structure
```
/project-root
├── /cmd                # Entry point for the application
│   └── main.go         # Main application file
├── /internal           # Internal packages
│   └── /handlers       # HTTP handlers
│       └── handlers.go
├── /web                # Web assets
│   ├── /templates      # HTML templates
│   │   └── index.html
│   └── /static         # Static files (CSS, JS, images)
│       ├── style.css
│       └── script.js
├── go.mod              # Go module file
└── README.md           # Project documentation
```

## Getting Started
1. **Install Go**: Ensure you have Go installed on your system.
2. **Clone the Repository**: Clone this project to your local machine.
3. **Run the Application**:
   - Navigate to the `cmd` directory.
   - Run `go run main.go` to start the server.

## Features
- **Go Backend**: Handles HTTP requests and serves HTML templates.
- **HTML Frontend**: Includes basic templates for rendering pages.
- **Static Assets**: CSS, JavaScript, and images are organized in the `/web/static` directory.

## Development
- Add new HTTP handlers in `/internal/handlers`.
- Create or update HTML templates in `/web/templates`.
- Add static files (CSS, JS, images) in `/web/static`.

## License
This project is licensed under the MIT License. See the LICENSE file for details.
