# Snippetbox

A web application for sharing text snippets, built with Go. Uses PostgreSQL for data persistence, server-side sessions, and Go's `html/template` package for rendering.

## Features

- Create, view, and list text snippets with configurable expiry (1 day, 7 days, or 1 year)
- Server-side session management with PostgreSQL-backed session store
- Flash messages for user feedback
- Form validation with error display
- Security headers (CSP, X-Frame-Options, etc.)
- Request logging and panic recovery middleware
- RESTful JSON API endpoint

## Project Structure

```
.
├── cmd/
│   └── web/                    # Web server application code
│       ├── handlers.go         # HTTP handler functions
│       ├── helpers.go          # Helper methods (error handling, template rendering)
│       ├── main.go             # Application entry point, server config
│       ├── middleware.go       # HTTP middleware (logging, security, panic recovery)
│       ├── routes.go           # Route definitions and middleware chains
│       └── templates.go        # Template cache for HTML rendering
├── internal/
│   ├── models/
│   │   ├── errors.go           # Custom error types (ErrNoRecord)
│   │   └── snippets.go         # Snippet database model and queries
│   └── validator/
│       └── validator.go        # Form validation helpers
├── ui/
│   ├── html/
│   │   ├── base.tmpl.html      # Base layout template
│   │   ├── pages/
│   │   │   ├── create.tmpl.html # Snippet creation form
│   │   │   ├── home.tmpl.html   # Home page template
│   │   │   └── view.tmpl.html   # Snippet view page template
│   │   └── partials/
│   │       └── nav.tmpl.html   # Navigation partial
│   └── static/                 # Static assets (CSS, JS, images)
├── tmp/                        # Temporary files (Air hot reload)
└── go.mod                      # Go module definition
```

## Routes

| Method | Path                | Handler             | Description                    |
| ------ | ------------------- | ------------------- | ------------------------------ |
| GET    | `/`                 | `home`              | Home page (lists snippets)     |
| GET    | `/snippet/view/:id` | `snippetView`       | View a specific snippet        |
| GET    | `/snippet/create`   | `snippetCreate`     | Show snippet creation form     |
| POST   | `/snippet/create`   | `snippetCreatePost` | Handle snippet form submission |
| GET    | `/snippets`         | `snippetLatest`     | Get latest snippets (JSON API) |
| GET    | `/static/*filepath` | `fileServer`        | Serve static assets            |

## Middleware

| Middleware      | Description                                        |
| --------------- | -------------------------------------------------- |
| `secureHeaders` | Sets security headers (CSP, X-Frame-Options, etc.) |
| `logRequest`    | Logs incoming HTTP requests                        |
| `recoverPanic`  | Recovers from panics and returns 500 error         |
| `LoadAndSave`   | Session management (load/save session data)        |

## Dependencies

| Package                    | Purpose                                |
| -------------------------- | -------------------------------------- |
| `jackc/pgx/v5`             | PostgreSQL driver & connection pool    |
| `alexedwards/scs/v2`       | Session management                     |
| `alexedwards/scs/pgxstore` | PostgreSQL session store               |
| `julienschmidt/httprouter` | HTTP request router                    |
| `justinas/alice`           | Middleware chaining                    |
| `go-playground/form`       | Form decoding                          |
| `joho/godotenv`            | Load environment variables from `.env` |
| `google/uuid`              | UUID generation for snippet IDs        |

## Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL

### Environment Variables

| Variable          | Description                  | Default          |
| ----------------- | ---------------------------- | ---------------- |
| `DATABASE_URL`    | PostgreSQL connection string | —                |
| `SNIPPETBOX_ADDR` | HTTP server address          | `127.0.0.1:4000` |

Create a `.env` file in the project root or export the variables directly.

### Running the Application

```bash
go run ./cmd/web
```

Flags:

- `-addr` — HTTP network address (overrides `SNIPPETBOX_ADDR`)
- `-dsn` — PostgreSQL connection string (overrides `DATABASE_URL`)

The server starts on `http://127.0.0.1:4000` by default.

### Development with Hot Reload

```bash
air
```
