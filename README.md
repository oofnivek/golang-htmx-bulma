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
    make run
    # or
    go run ./cmd/web
    ```
4.  **Generating Password Hashes (for Database Seeding)**:
    To manually insert or update a user in your MySQL database with a compatible hashed password, you can use the built-in test helper:
    - Open `internal/pkg/crypto/password_test.go` and set your desired plain password.
    - Run the generator command:
      ```bash
      go test -v ./internal/pkg/crypto -run TestGenerateHash
      ```
    - Use the printed Base64 hash in your SQL insert or update query.


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

## Available Commands

| Command | Description |
|---|---|
| `make run` | Run the application |
| `make build` | Compile the binary |
| `make test` | Run all tests |
| `make test-coverage` | Run tests with filtered coverage report |
| `make vet` | Run `go vet` |
| `make lint` | Run `golangci-lint` |
| `make clean` | Remove coverage artifacts |

## Technology Stack

- **Backend**: [Gin](https://github.com/gin-gonic/gin)
- **Frontend**: [HTMX](https://htmx.org/)
- **CSS**: [Bulma](https://bulma.io/)
- **Database**: MySQL
