# Snippetbox

A web application for sharing text snippets, built with Go.

## Project Structure

```
.
├── cmd/
│   └── web/         # Web server application code
│       ├── handlers.go    # HTTP handler functions
│       └── main.go        # Application entry point, routing
├── internal/              # Private application code (not importable) - e.g. database models, utilities
├── ui/
│   ├── html/              # HTML templates
│   └── static/            # Static assets (CSS, JS, images)
├── tmp/                   # Temporary files (Air hot reload)
├── .air.toml              # Air configuration for live reloading
└── go.mod                 # Go module definition
```

## Directory Descriptions

| Directory   | Purpose                                                    |
| ----------- | ---------------------------------------------------------- |
| `cmd/web`   | Application-specific code for the web server               |
| `internal`  | Private packages that can only be imported by this project |
| `ui/html`   | HTML templates for rendering pages                         |
| `ui/static` | Static files served directly (CSS, JavaScript, images)     |
| `tmp`       | Temporary build output from Air hot reloader               |
| `.air.toml` | Air configuration for live reloading                        |
| `go.mod`    | Go module definition                                        |

## Routes

| Method | Path                | Handler         | Description             |
| ------ | ------------------- | --------------- | ----------------------- |
| GET    | `/`                 | `home`          | Home page               |
| GET    | `/snippet/view?id=` | `snippetView`   | View a specific snippet |
| POST   | `/snippet/create`   | `snippetCreate` | Create a new snippet    |

## Getting Started

### Prerequisites

- Go 1.25+

### Running the Application

```bash
go run ./cmd/web
```

The server starts on `http://127.0.0.1:4000`

### Development with Hot Reload

```bash
air
```
