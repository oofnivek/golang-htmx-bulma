# Go + HTMX + Bulma Project

A full-stack web application built with Go, Gin, HTMX, Bulma CSS, and MySQL.

## Project Structure

- `cmd/web/`: Main entrypoint
- `internal/`: Private application domain logic
  - `internal/user/`: Core User, Role, and Auth domain logic (interfaces, local DB implementations, HTTP API clients)
  - `internal/vehicle/`: Core Vehicle (Colors & Makes) domain logic
  - `internal/http/`: Web & REST handlers, middleware, and route registration
  - `internal/view/`: HTML render helpers
  - `internal/config/`: Configuration manager
  - `internal/db/`: Database helpers
- `templates/`: HTML templates (layouts, pages, partials)
- `static/`: Static assets (CSS, JS, images)
- `migrations/`: MySQL schema migrations

## Getting Started

1.  **Prerequisites**: Go 1.26+, MySQL
2.  **Setup Environment**:
    ```bash
    cp .env.example .env
    # Update .env with your MySQL DSN settings and JWT signature key
    ```

3.  **Run the Application (Multi-Role Configuration)**:

    The application is built as a modular monolith and can be run in three distinct roles using the `APP_ROLE` environment variable:

    ### Option A: Monolith Mode (Default)
    Runs all web presentation handlers and API endpoints in-memory, requiring connections to both User and Vehicle MySQL databases.
    ```bash
    go run ./cmd/web
    # Server starts on port 8080 (or PORT in .env)
    ```

    ### Option B: Isolated User Service Mode
    Runs strictly as a headless JSON API server on port `8081`, requiring connection **only** to the user database.
    ```bash
    APP_ROLE=user-service PORT=8081 go run ./cmd/web
    ```

    ### Option C: Isolated Web/View Mode
    Runs as a pure presentation layer on port `8080`, serving HTML/HTMX templates and connecting only to the vehicle database. User, Role, and Auth actions are dynamically delegated to the remote User Service over REST APIs.
    ```bash
    APP_ROLE=web-view PORT=8080 USER_SERVICE_URL="http://localhost:8081" go run ./cmd/web
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
| `internal/user/` | **High (80%+)** (Business logic, REST clients, and DB repositories) |
| `internal/vehicle/` | **High (80%+)** (Business logic and DB repositories) |
| `internal/http/handlers/` | Medium |
| `internal/view/`, `internal/config/`, `internal/db/`, `cmd/` | Skip (infrastructure/wiring) |

## Available Commands

| Command | Description |
|---|---|
| `make run` | Run the application (Monolith mode) |
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
