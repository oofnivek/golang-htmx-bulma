# Go + HTMX + Bulma Project

A full-stack web application built with Go, Gin, HTMX, Bulma CSS, and MySQL.

## Project Structure

- `cmd/web/`: Main entrypoint
- `internal/`: Application logic (config, db, handlers, routes, etc.)
- `templates/`: HTML templates (layouts, pages, partials)
- `static/`: Static assets (CSS, JS, images)
- `migrations/`: MySQL schema migrations

## Getting Started

1.  **Prerequisites**: Go 1.26+, MySQL
2.  **Setup Environment**:
    ```bash
    cp .env.example .env
    # Update .env with your MySQL DSN
    ```
3.  **Run the Application**:
    ```bash
    go run ./cmd/web
    ```

## Testing & Coverage

Run all tests:
```bash
make test
# or
go test ./...
```

Run tests with a **filtered coverage report** (excludes infrastructure/wiring packages like `cmd/`, `internal/view/`, `internal/config/`, `internal/db/`, `internal/http/routes/`):
```bash
make test-coverage
```

Run tests with the **full coverage report** (all packages):
```bash
go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out
```

Open the **interactive HTML coverage report** in your browser:
```bash
go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
```

Clean up coverage artifacts:
```bash
make clean
```

### Coverage priorities

| Package | Target |
|---|---|
| `internal/service/` | 80%+ |
| `internal/repository/` | 80%+ |
| `internal/http/handlers/web/` | Medium |
| `internal/view/`, `internal/config/`, `internal/db/`, `cmd/` | Skip (infrastructure/wiring) |

## Technology Stack

- **Backend**: [Gin](https://github.com/gin-gonic/gin)
- **Frontend**: [HTMX](https://htmx.org/)
- **CSS**: [Bulma](https://bulma.io/)
- **Database**: MySQL
