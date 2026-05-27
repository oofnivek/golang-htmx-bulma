# AGENTS.md

## Project overview
This project is a full-stack web application built with Go, Gin, HTMX, and MySQL, using a custom Claude design template for styling.

- Backend framework: Gin
- Frontend: server-rendered HTML templates enhanced with HTMX
- CSS framework: Claude design template (custom)
- Database: MySQL only (connects to existing databases — no migrations run from this repo)
- Goal: keep the app simple, server-driven, and easy to maintain

## Database schema
DDL for all tables is in `schema/`, organized by database name:
- `schema/vehicles/` — tables in the `vehicles` database
- `schema/fms_users/` — tables in the `fms_users` database

## Core principles
- Prefer server-rendered HTML over heavy frontend JavaScript.
- Use HTMX for progressive enhancement, partial updates, and simple interactivity.
- Use Claude design template components and classes before adding custom CSS.
- Keep dependencies small unless there is a clear reason.
- Prefer clarity and maintainability over clever abstractions.
- Make small, focused changes instead of large rewrites.

## New feature scaffold

Use this checklist when adding a new CRUD feature. Work in layer order — each layer's interface is the contract the next layer depends on. Replace `Thing` / `thing` / `things` with the real name.

### Layer 1 — Domain model  
File: `internal/<domain>/<thing>.go`
```go
type Thing struct {
    ID        int64      `json:"id"`
    Name      string     `json:"name"`
    // ... domain fields
    CreatedBy string     `json:"created_by"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedBy *string    `json:"updated_by"`
    UpdatedAt *time.Time `json:"updated_at"`
}
```
- Use `int64` for numeric PKs; `string` for natural-key PKs (e.g. role ID).
- Pointer types (`*string`, `*time.Time`) for nullable columns.
- Include display-only join fields (e.g. `TypeName string`) for list views.

### Layer 2 — Repository  
File: `internal/<domain>/<thing>_repository.go`

Interface (define at top of file, implemented below):
```go
type ThingRepository interface {
    GetAll() ([]Thing, error)
    GetPaged(limit, offset int, sortBy, sortOrder string) ([]Thing, error)
    Count() (int, error)
    GetByID(id int64) (*Thing, error)
    Create(t *Thing) error
    Update(t *Thing) error
    Delete(id int64) error
}
```
Implementation struct: `mysqlThingRepository` wrapping `*sql.DB`.  
Constructor: `NewThingRepository(db *sql.DB) ThingRepository`.  
Use named `const` for SELECT columns and FROM/JOIN clauses; extract a `scanThing()` helper.  
`GetPaged`: apply a `sortableColumns` allowlist map before interpolating `sortBy`/`sortOrder` into the query.

### Layer 3 — Service  
File: `internal/<domain>/<thing>_service.go`

Interface:
```go
type ThingService interface {
    ListAll() ([]Thing, error)
    ListPaged(page, pageSize int, sortBy, sortOrder string) ([]Thing, int, error)
    FindByID(id int64) (*Thing, error)
    CreateThing(/* fields */, user string) (*Thing, error)
    UpdateThing(id int64, /* fields */, user string) (*Thing, error)
    DeleteThing(id int64) error
}
```
Implementation struct: `thingService{repo ThingRepository}`.  
Constructor: `NewThingService(repo ThingRepository) ThingService`.  
`ListPaged`: clamp `page`/`pageSize` to ≥1, compute `offset = (page-1)*pageSize`, call `repo.GetPaged` then `repo.Count`.  
`CreateThing`: set `CreatedAt`, `UpdatedAt` to `time.Now().UTC()`; set `CreatedBy`, `UpdatedBy` from `user`.  
`UpdateThing`: fetch existing record first (`GetByID`), return `nil,nil` if not found, mutate then call `repo.Update`.

### Layer 4 — Remote service (for web-view mode)  
File: `internal/<domain>/remote_service.go` — add a new `remoteThingService` struct at the bottom.

```go
type remoteThingService struct{ baseURL string }
func NewRemoteThingService(baseURL string) ThingService { ... }
```
Each method: build URL, call `http.Get`/`http.Post`/etc., decode JSON response.  
Mirror the local service's method signatures exactly.

### Layer 5 — Web handler  
File: `internal/http/handlers/web/<thing>.go`

```go
type ThingHandler struct{ svc domain.ThingService }
func NewThingHandler(svc domain.ThingService) *ThingHandler
```
Methods: `List`, `View`, `CreateForm`, `Create`, `EditForm`, `Update`, `Delete`, `DeleteConfirm`.  
- `List`: read `page`, `pageSize`, `sortBy`, `sortOrder`, `tz` from query; render `pages/things/index.html` with pagination map.  
- `CreateForm` / `EditForm`: render `pages/things/form.html`; pass `"action"` URL and any lookup lists needed by dropdowns.  
- `Create` / `Update`: parse form values, call service, `c.Redirect(303, "/things")`.  
- `Delete`: call service, `c.Status(200)` (HTMX row removal via `hx-target`/`hx-swap="outerHTML"`).  
- `DeleteConfirm`: fetch record, render `partials/modals/delete_confirm.html` with `Name`, `DeleteURL`, `RowID`.  
- `View`: fetch record, render `pages/things/view.html`.

### Layer 6 — API handler  
File: `internal/http/handlers/api/<thing>.go`

```go
type ThingAPIHandler struct{ svc domain.ThingService }
func NewThingAPIHandler(svc domain.ThingService) *ThingAPIHandler
```
Methods: `ListAll` (no pagination), `List` (paged), `Get`, `Create`, `Update`, `Delete`.  
Use `c.ShouldBindJSON` for request bodies; return `c.JSON(400/404/500/200/201, gin.H{...})`.

### Layer 7 — Routes  
File: `internal/http/routes/routes.go`

1. Add `thingHandler *web.ThingHandler` and `thingAPI *api.ThingAPIHandler` parameters to `RegisterRoutes`.
2. Add to the nil-guard condition (`if ... || thingHandler != nil ...`).
3. Inside `protected` group:
```go
if thingHandler != nil {
    g := protected.Group("/things")
    {
        g.GET("",          thingHandler.List)
        g.GET("/new",      thingHandler.CreateForm)
        g.POST("",         thingHandler.Create)
        g.GET("/:id/view", thingHandler.View)
        g.GET("/:id/edit", thingHandler.EditForm)
        g.POST("/:id",     thingHandler.Update)
        g.DELETE("/:id",   thingHandler.Delete)
        g.GET("/:id/delete", thingHandler.DeleteConfirm)
    }
}
```
4. Inside `apiGroup`:
```go
if thingAPI != nil {
    apiGroup.GET("/things/all", thingAPI.ListAll)
    apiGroup.GET("/things",     thingAPI.List)
    apiGroup.GET("/things/:id", thingAPI.Get)
    apiGroup.POST("/things",    thingAPI.Create)
    apiGroup.PUT("/things/:id", thingAPI.Update)
    apiGroup.DELETE("/things/:id", thingAPI.Delete)
}
```

### Layer 8 — Main wiring  
File: `cmd/web/main.go`

Add variable declarations, then wire inside **each APP_ROLE case** that owns the domain:
```go
// vehicle-service / user-service case:
thingRepo := domain.NewThingRepository(db)
thingSvc  := domain.NewThingService(thingRepo)
thingHandler = web.NewThingHandler(thingSvc)
thingAPI     = api.NewThingAPIHandler(thingSvc)

// web-view case (remote):
thingSvc = domain.NewRemoteThingService(serviceURL)
thingHandler = web.NewThingHandler(thingSvc)
```
Pass `thingHandler` and `thingAPI` to `routes.RegisterRoutes(...)`.

### Layer 9 — Templates

| File | Purpose |
|---|---|
| `templates/pages/things/index.html` | Paginated table, sort headers, rows per page, timezone selector, pagination nav |
| `templates/pages/things/form.html` | Create/edit form; `autofocus` on first field; Save + Cancel footer buttons |
| `templates/pages/things/view.html` | Read-only mirror of form; only "Back to List" footer button |
| `templates/partials/things/table_row.html` | `<tr id="thing-row-{{ .ID }}">` with icon-only action buttons |

Template conventions recap:
- Index: `hx-target="body" hx-push-url="true"` on sort links, page links, and selectors.
- Table row: `{{ define "things/table_row.html" }}` — called with `{{ template "things/table_row.html" . }}`.
- Action button order: all three buttons — view (`ra-btn`, eye), edit (`ra-btn`, pencil), delete (`ra-btn danger`, trash) — must always be present in the Actions column, in that order. Wrap in a `row-act` div.
- Form footer: Save button (`btn btn-primary`, `fa-floppy-disk`) **left of** Cancel button (`btn`). Wrap in a `form-footer` div.
- Delete modal target: `hx-target="#modal-container"`.

### Layer 10 — Sidebar nav  
File: `templates/layouts/base.html` — add an `<a href="/things" class="sb-item">` link with the appropriate FontAwesome icon inside a `<span class="sb-ico">` under the correct `<div class="sb-section">` section.

---

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
## Current Clean Modular Monolith Architecture

The application can run in three distinct modes controlled by the `APP_ROLE` environment variable:

- **web-view** – Starts the Gin server and serves HTML pages. It uses the remote **user-service** and **vehicle-service** via HTTP (`USER_SERVICE_URL`, `VEHICLE_SERVICE_URL`).
- **user-service** – Runs only the user domain (handlers, services, repositories) on its own port, exposing both web and API endpoints for user operations.
- **vehicle-service** – Runs only the vehicle domain (handlers, services, repositories) on its own port, exposing both web and API endpoints for vehicle operations.

Each role wires only the required internal packages, keeping the other domain packages compiled but unused. The Mermaid diagram in the project README visualises this architecture and the wiring between `cmd/web/main.go`, the route registration, and the service layers.


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
  - **Searching**: Searching must NOT be auto-triggered on every keyup. It must only be triggered when the user clicks the search button. Searching via the Enter key must be disabled. The search button must be disabled when the search term is exactly 1 character long, and enabled when it is 0 characters (to reset) or 2+ characters long.

## Auth conventions
- Auth is enforced through Gin middleware and route groups, not by folder names.
- **JWT Implementation**:
  - Use `github.com/golang-jwt/jwt/v5` for token management.
  - Algorithm: **HS256**.
  - Signing Key: Minimum **256 bits (32 bytes)** random key, stored as a 64-character hex string in `.env`.
  - Claims: Standard claims (`sub`, `iss`, `iat`, `exp`, `jti`) plus custom claims (`email`, `role_id`, `first_name`, `last_name`).
  - Storage: Store the JWT in an **HttpOnly** cookie named `jwt_token`.
- **Middleware**:
  - Protected routes must be wrapped in an `AuthMiddleware` group.
  - The middleware must handle both standard redirects and **HTMX redirects** (using the `HX-Redirect` header) to ensure a seamless UI experience when a session expires.
- Protected browser routes and protected HTMX routes must pass auth checks.
- Protected API routes must pass auth checks.
- Keep web handlers and API handlers separate when that improves clarity.
- Do not assume an HTMX request is trusted just because it comes from the frontend.
- If bearer token auth is used for a route, validate it in middleware or a well-defined auth layer.
- If a future API is intended for service-to-service use, define its auth rules explicitly in routes and middleware.

## Technology rules
- Prefer the Go standard library unless an external package clearly improves maintainability.
- Use the current project templating approach consistently.
- Use Claude design template classes consistently for layout, forms, buttons, tables, notifications, and modals.
- Add custom JavaScript only when HTMX cannot reasonably solve the interaction.
- Add custom CSS only when the Claude design template cannot reasonably handle the styling.

## Database rules
- MySQL is the only supported database.
- All schema design, queries, migrations, and repository code should target MySQL.
- Use MySQL-compatible SQL syntax and features.
- **Search Inefficiency**: Avoid using `LIKE '%keyword%'` for searching as it performs a full table scan and is extremely inefficient on large datasets. Use `MATCH() AGAINST()` with a `FULLTEXT` index (preferably with the `ngram` parser for partial matches) instead.
- **Timestamps**: Always store timestamps in UTC in the database.
- **Row Creation**: Upon record creation, set `created_at` and `updated_at` to the same timestamp, and set `created_by` and `updated_by` to the authenticated user's email (`user_email` from Gin context).
- **Record Update**: When editing a record, only `updated_at` and `updated_by` are changed; `created_at` and `created_by` remain unchanged.
- Keep SQL organized in repository/store layers.
- Document any MySQL-specific assumptions when they affect schema, indexing, transactions, or query behavior.
- Avoid raw SQL duplication across multiple files when a shared repository method is more maintainable.
- Ask before making destructive schema changes such as dropping columns, renaming tables, or irreversible migrations.

## Configuration rules
- All environment-specific behavior should be controlled through config.
- Do not hardcode DSN, ports, secrets, or file paths.
- Prefer environment variables or config files for:
  - VEHICLE_DB_DSN
  - FMS_USER_DB_DSN
  - JWT_SIGNING_KEY
  - server port
  - app environment
  - auth settings
  - session/security settings

## Logging and Telemetry rules
- **OpenTelemetry Standard**: All application logging must use OpenTelemetry standard logging via the OTel slog bridge.
- **Structured JSON Logging**: Do not use `log.Printf`, `log.Println`, or `fmt.Printf` for application logging. All logging must use standard Go's structured logger `log/slog`.
- **Vendor Neutrality**: Keep the codebase free from proprietary vendor SDKs (e.g. Datadog SDKs). Use OTel standard exporters to ship logs to OTel-compliant backends.
- **Rich Context/Attributes**: Always include useful structured attributes as key-value pairs (e.g., `email`, `client_ip`, `error`) so they are indexed automatically as searchable facets in log backends.
- **Secure Server Logs**: Use `slog.Error` to log actual database/system failures securely on the server-side, while returning safe, user-friendly error messages to the client. Never pass `err.Error()` directly to `c.String(500, ...)` — that leaks internal details to the browser.
- **Handler error logging**: Every handler path that returns a 5xx must call `slog.Error` before returning. Always include the operation context and the structured `"error"` attribute:
  ```go
  slog.Error("failed to delete vehicle color", "id", id, "error", err)
  c.String(http.StatusInternalServerError, "Failed to delete vehicle color")
  ```
  Include the record ID (or other identifying field) so log entries are actionable without a database query. Use a consistent message format: `"failed to <verb> <entity>"` (e.g. `"failed to list vehicle makes"`, `"failed to create fuel type"`).

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

## Package design rules and Modular Monolith Boundaries
- **Encapsulated Domain Packages**: Group related domain features (e.g., models, repositories, and services) into self-contained packages under `internal/` (such as `internal/vehicle`).
- **Strict Boundary Check / Warnings**:
  - The AI MUST actively warn the user if they attempt to put new or existing domain files in the wrong service package, or if they import/leak private package details across boundaries.
  - The AI MUST verify imports and refuse or warn when domain-specific models/logic bypass the public interface of a domain package.
- Do not create deep package nesting without a clear reason.
- Do not create generic utility packages too early.
- Prefer domain- or feature-oriented packages over vague shared helpers when the project grows.
- Keep package APIs minimal and explicit.
- Do not introduce `pkg/` unless code is intentionally meant for reuse outside this module.

## Template conventions
- Keep templates organized into layouts, pages, and partials.
- Reuse partials for repeated UI such as forms, flash messages, table rows, and modals.
- **Autofocus**: When creating "new", "edit", or "delete" (modal) views, always add the `autofocus` attribute to the first or most relevant input field to improve UX.
- **Cancel Button**: Form pages ("new" and "edit") must include a "Cancel" button in the `form-footer` area that redirects back to the index page using HTMX. The "Save" button must be on the left of the "Cancel" button, use the `btn btn-primary` class, and include a FontAwesome `fa-floppy-disk` icon. The "Cancel" button must use the `btn` class.
- **View Page**: Modules with extensive details must provide a read-only view page. This page should mirror the edit form's layout but with all inputs in `readonly` or `disabled` mode. The footer must only contain a "Back to List" button (using `btn` and right-aligned).
- **Action Buttons**: The "Actions" column in index tables must use icon-only buttons (FontAwesome) rather than text. Use `ra-btn` for view (eye icon) and edit (pencil icon), and `ra-btn danger` for delete (trash icon). Wrap in a `row-act` div.
- **Timezone selector**: Use the shared `web.Timezones` constant when populating timezone dropdowns in index pages (e.g., vehicle_colors, vehicle_makes, users) to keep the list consistent across the application.
- Keep conditional logic in templates minimal.
- Format HTML clearly so partials are easy to update.
- Keep HTMX-targeted fragments small and reusable.

## Claude design template conventions
- Use the custom design template's component classes (`btn`, `btn-primary`, `card`, `tbl`, `form-card`, `form-field`, `form-label`, `form-input`, `form-footer`, `ph`, `ra-btn`, `row-act`, `pill`, etc.) consistently.
- Keep styling consistent with the design template's CSS custom properties (`var(--bg)`, `var(--ink-2)`, `var(--ink-3)`, etc.).
- Avoid large custom CSS overrides that conflict with the design template.
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

## Legacy / ignored fields

- `VehicleType.OldID` (`old_id` column) — legacy migration artifact. The column exists in the database and on the struct but must not be exposed in forms, API responses, or new features.

## Ask first
Ask before doing any of the following:
- adding or removing dependencies
- changing schema in a destructive way
- changing the folder structure significantly
- replacing Gin, HTMX, the Claude design template, or MySQL
- introducing a client-side frontend framework
- modifying authentication or authorization behavior significantly

## Preferred workflow for changes
When implementing a feature or fix:
1. Understand the existing route, handler, service, repository, and template flow.
2. Change the smallest number of files needed.
3. Keep full-page routes, HTMX partial routes, and API routes clearly separated when practical.
4. Reuse existing patterns before introducing new ones.
5. Update or add tests when behavior changes.
6. Ensure the result matches the Go + Gin + HTMX + Claude design template + MySQL style used in this repo.

## Good defaults for this project
- Gin for routing and middleware
- Server-driven UI
- Thin handlers
- Clear services
- Repository-based data access
- Reusable HTML partials
- Minimal JavaScript
- Claude design template styling
- MySQL-only data model (this app connects to an existing database — no migrations are run from this repo)
- `cmd/` for entrypoints and `internal/` for app code
- **Time Handling**: UTC in database, local/selected timezone in UI