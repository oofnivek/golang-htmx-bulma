# AGENTS.md

## Project overview
This project is a full-stack web application built with Go, Gin, HTMX, Bulma CSS, and MySQL.

- Backend framework: Gin
- Frontend: server-rendered HTML templates enhanced with HTMX
- CSS framework: Bulma
- Database: MySQL only
- Goal: keep the app simple, server-driven, and easy to maintain

## Core principles
- Prefer server-rendered HTML over heavy frontend JavaScript.
- Use HTMX for progressive enhancement, partial updates, and simple interactivity.
- Use Bulma components and classes before adding custom CSS.
- Keep dependencies small unless there is a clear reason.
- Prefer clarity and maintainability over clever abstractions.
- Make small, focused changes instead of large rewrites.

## Go project structure
Use idiomatic Go project structure with `cmd/` for the executable entrypoint and `internal/` for application-specific code that should not be imported by other modules.

Preferred layout:

```text
.
├── cmd/
│   └── web/
│       └── main.go
├── internal/
│   ├── config/
│   ├── db/
│   ├── model/
│   ├── repository/
│   ├── service/
│   ├── view/
│   └── http/
│       ├── handlers/
│       │   ├── web/
│       │   └── api/
│       ├── middleware/
│       └── routes/
├── templates/
│   ├── layouts/
│   ├── pages/
│   └── partials/
├── static/
│   ├── css/
│   ├── js/
│   └── img/
├── migrations/
├── .env.example
├── go.mod
├── go.sum
├── README.md
└── AGENTS.md
```

Structure rules:
- Put the application entrypoint in `cmd/web/main.go`.
- Keep `main.go` focused on startup and dependency wiring.
- Put application-specific Go code under `internal/`.
- In this project, the HTTP layer uses `handlers` rather than `controllers`.
- Put Gin web handlers in `internal/http/handlers/web/`.
- Put Gin API handlers in `internal/http/handlers/api/`.
- Put Gin middleware in `internal/http/middleware/`.
- Put route registration in `internal/http/routes/`.
- Put business logic in `internal/service/`.
- Put MySQL access code in `internal/repository/`.
- Put database connection setup and database helpers in `internal/db/`.
- Put models or domain structs in `internal/model/`.
- Put view models or render helpers in `internal/view/` if needed.
- Put shared page layouts in `templates/layouts/`.
- Put full page templates in `templates/pages/`.
- Put reusable partials and HTMX fragments in `templates/partials/`.
- Put CSS, JavaScript, and images in `static/`.
- Put schema migrations in `migrations/`.
- Keep the structure simple; only add new packages when there is clear complexity that justifies them.

## Why `internal/` is used
- `internal/` is used to mark application code as private to this module.
- Packages under `internal/` are implementation details, not public library APIs.
- Authentication rules are not controlled by `internal/`; authentication is controlled by Gin middleware and route groups.
- Do not use `internal/` to separate authenticated routes from unauthenticated routes.

## Architecture rules
- Gin handlers should stay thin; handlers should parse input, call services, and return responses.
- Put business logic in services/use-cases, not in Gin handlers or templates.
- Put database access in repositories.
- Templates should contain presentation logic only, not business logic.
- Do not put SQL in handlers or templates.
- Do not mix MySQL access directly into handlers.
- Reuse existing patterns before introducing new abstractions.

## Gin and routing conventions
- Use Gin routing, middleware, and context patterns consistently.
- Prefer route groups for related features.
- Use route groups to apply auth, logging, rate limiting, or other shared middleware.
- Use `c.HTML(...)` for server-rendered pages and HTML fragments.
- Use `c.JSON(...)` only for API-style endpoints or when HTML is not appropriate.
- Use `c.Param(...)`, `c.Query(...)`, and form binding carefully and explicitly.
- Return proper HTTP status codes.
- Keep middleware focused and reusable.
- Do not place business logic inside middleware.

## HTTP and HTMX conventions
- Gin handlers serve both normal browser requests and HTMX requests.
- Standard page routes should render full pages.
- HTMX routes should return partial HTML snippets when appropriate.
- If the same route supports both full-page and HTMX requests, detect the HTMX request and choose the correct rendering mode.
- For browser interactions, prefer returning HTML fragments from Gin handlers instead of JSON when the client is HTMX.
- Preserve progressive enhancement where practical so core flows still work without HTMX.
- For validation errors, return user-friendly HTML fragments or page responses suitable for the request type.
- **Index Pages**: 
  - Results must always be paginated.
  - The number of rows per page must be customizable by the user.
  - Every column must be sortable (ascending/descending), except for the "Actions" column.
  - Sorting and pagination state (including timezone) must be preserved across interactions using HTMX.
  - **Searching**: Searching must NOT be auto-triggered on every keyup. It must only be triggered when the user clicks the search button or presses the Enter key.

## Auth conventions
- Auth is enforced through Gin middleware and route groups, not by folder names.
- Protected browser routes and protected HTMX routes must pass auth checks.
- Protected API routes must pass auth checks.
- Keep web handlers and API handlers separate when that improves clarity.
- Do not assume an HTMX request is trusted just because it comes from the frontend.
- If bearer token auth is used for a route, validate it in middleware or a well-defined auth layer.
- If a future API is intended for service-to-service use, define its auth rules explicitly in routes and middleware.

## Technology rules
- Prefer the Go standard library unless an external package clearly improves maintainability.
- Use the current project templating approach consistently.
- Use Bulma classes consistently for layout, forms, buttons, tables, notifications, and modals.
- Add custom JavaScript only when HTMX cannot reasonably solve the interaction.
- Add custom CSS only when Bulma cannot reasonably handle the styling.

## Database rules
- MySQL is the only supported database.
- All schema design, queries, migrations, and repository code should target MySQL.
- Use MySQL-compatible SQL syntax and features.
- **Search Inefficiency**: Avoid using `LIKE '%keyword%'` for searching as it performs a full table scan and is extremely inefficient on large datasets. Use `MATCH() AGAINST()` with a `FULLTEXT` index (preferably with the `ngram` parser for partial matches) instead.
- **Timestamps**: Always store timestamps in UTC in the database.
- **Row Creation**: Upon row creation, set `created_at` and `updated_at` to the same value, and set `created_by` and `updated_by` to the same value.
- Keep SQL organized in repository/store layers.
- Document any MySQL-specific assumptions when they affect schema, indexing, transactions, or query behavior.
- Avoid raw SQL duplication across multiple files when a shared repository method is more maintainable.
- Ask before making destructive schema changes such as dropping columns, renaming tables, or irreversible migrations.

## Configuration rules
- All environment-specific behavior should be controlled through config.
- Do not hardcode DSN, ports, secrets, or file paths.
- Prefer environment variables or config files for:
  - VEHICLE_DB_DSN
  - server port
  - app environment
  - auth settings
  - session/security settings

## Code style
- Write idiomatic Go.
- **Language**: Always use US English for variable names, database fields, and UI text (e.g., "color" instead of "colour").
- Keep packages focused and cohesive.
- Keep functions small and focused.
- Prefer descriptive names over short clever names.
- Use short, lowercase package names.
- Keep interfaces small and define them where they are used.
- Return errors explicitly; do not ignore them.
- Wrap errors with useful context.
- Avoid global mutable state unless clearly justified.
- Prefer composition over unnecessary abstraction.

## Package design rules
- Do not create deep package nesting without a clear reason.
- Do not create generic utility packages too early.
- Prefer domain- or feature-oriented packages over vague shared helpers when the project grows.
- Keep package APIs minimal and explicit.
- Do not introduce `pkg/` unless code is intentionally meant for reuse outside this module.

## Template conventions
- Keep templates organized into layouts, pages, and partials.
- Reuse partials for repeated UI such as forms, flash messages, table rows, and modals.
- **Autofocus**: When creating "new", "edit", or "delete" (modal) views, always add the `autofocus` attribute to the first or most relevant input field to improve UX.
- Keep conditional logic in templates minimal.
- Format HTML clearly so partials are easy to update.
- Keep HTMX-targeted fragments small and reusable.

## Bulma conventions
- Prefer Bulma-native components and naming.
- Keep styling consistent with Bulma defaults unless the project defines a custom design system.
- Avoid large custom CSS overrides that fight Bulma.
- If custom CSS is necessary, keep it small, scoped, and documented.

## Commands
Use the actual repo commands if they exist. Typical commands may include:

- Run app: `make run` or `go run ./cmd/web`
- Test: `make test` or `go test ./...`
- Test with filtered coverage report: `make test-coverage`
- Vet: `make vet` or `go vet ./...`
- Lint: `make lint` or `golangci-lint run`
- Build: `make build` or `go build ./cmd/web`
- Clean coverage artifacts: `make clean`

If a command is unavailable in this repo, do not invent it; use the closest existing project command.

## Testing and validation
Before finishing changes, always do the relevant checks:

- Run tests using `go test ./...`.
- Keep `_test.go` files in the same directory/package as the code being tested.
- **Mocks**: Define mock implementations in the same package that uses them, and ensure the filename ends with `_test.go` (e.g., `mock_vehicle_color_repository_test.go`). This prevents them from being compiled into the production binary and automatically excludes them from code coverage reports.
- **Coverage priority**:
  - `internal/service/` — **High (80%+)**. Business logic and error paths must be tested.
  - `internal/repository/` — **High (80%+)**. SQL correctness and error handling must be tested.
  - `internal/http/handlers/` — **Medium**. Input parsing, response codes, and routing.
  - `internal/view/`, `internal/config/`, `internal/db/`, `cmd/` — **Low / Skip**. These are infrastructure wiring and bootstrap code. They are filesystem- or network-dependent, and failures are immediately visible at startup. Do not force unit tests onto them.
- Use **Mock Repositories** to test services in isolation without a real database.
- Run linting with `golangci-lint run` if available.
- Verify the app builds.
- Verify affected flows for both normal full-page requests and HTMX requests.
- Verify protected routes still enforce auth correctly.
- Verify MySQL-related queries, schema assumptions, and migrations still make sense.
- Prefer small, targeted tests for changed packages.
- **Do not automatically fix unit tests**: When implementing changes or new functionality, do not automatically update or fix unit tests unless explicitly asked by the USER. The USER may be iterating on functionality or UI and will request test and coverage fixes when they are done.

## Safety rules
- Never commit secrets, API keys, `.env` contents, or real credentials.
- Never delete large sections of code without clear reason.
- Never add heavy dependencies without justification.
- Ask before making destructive schema changes or irreversible migrations.
- Ask before changing core project structure or replacing major libraries/frameworks.

## Ask first
Ask before doing any of the following:
- adding or removing dependencies
- changing schema in a destructive way
- changing the folder structure significantly
- replacing Gin, HTMX, Bulma, or MySQL
- introducing a client-side frontend framework
- modifying authentication or authorization behavior significantly

## Preferred workflow for changes
When implementing a feature or fix:
1. Understand the existing route, handler, service, repository, and template flow.
2. Change the smallest number of files needed.
3. Keep full-page routes, HTMX partial routes, and API routes clearly separated when practical.
4. Reuse existing patterns before introducing new ones.
5. Update or add tests when behavior changes.
6. Ensure the result matches the Go + Gin + HTMX + Bulma + MySQL style used in this repo.

## Good defaults for this project
- Gin for routing and middleware
- Server-driven UI
- Thin handlers
- Clear services
- Repository-based data access
- Reusable HTML partials
- Minimal JavaScript
- Bulma-first styling
- MySQL-only data model and migrations
- `cmd/` for entrypoints and `internal/` for app code
- **Time Handling**: UTC in database, local/selected timezone in UI