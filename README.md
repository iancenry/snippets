# Snippetbox

A web application for sharing text snippets, built with Go. Uses PostgreSQL for data persistence and Go's `html/template` package for server-side rendering.

## Project Structure

```
.
├── cmd/
│   └── web/                    # Web server application code
│       ├── handlers.go         # HTTP handler functions
│       ├── helpers.go          # Helper methods (error handling)
│       ├── main.go             # Application entry point, server config
│       └── routes.go           # Route definitions and middleware
├── internal/
│   └── models/
│       └── snippets.go         # Snippet database model and queries
├── ui/
│   ├── html/
│   │   ├── base.tmpl.html      # Base layout template
│   │   ├── pages/
│   │   │   └── home.tmpl.html  # Home page template
│   │   └── partials/
│   │       └── nav.tmpl.html   # Navigation partial
│   └── static/                 # Static assets (CSS, JS, images)
├── tmp/                        # Temporary files (Air hot reload)
└── go.mod                      # Go module definition
```

## Routes

| Method | Path              | Handler         | Description             |
| ------ | ----------------- | --------------- | ----------------------- |
| GET    | `/`               | `home`          | Home page               |
| GET    | `/snippet/view/`  | `snippetView`   | View a specific snippet |
| POST   | `/snippet/create` | `snippetCreate` | Create a new snippet    |
| GET    | `/static/`        | `fileServer`    | Serve static assets     |

## Dependencies

| Package         | Purpose                                |
| --------------- | -------------------------------------- |
| `jackc/pgx/v5`  | PostgreSQL driver & connection pool    |
| `joho/godotenv` | Load environment variables from `.env` |
| `google/uuid`   | UUID generation for snippet IDs        |

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
