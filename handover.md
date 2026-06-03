# Handover — condo_car_park feature

## Status: Layers 1–8 complete, build passing. Layers 9–10 not yet started.

---

## What was done

### Schema
- Extracted `condo_car_park` DDL from the `vehicles` database into `schema/vehicles/condo_car_park.sql`.
- The table is a lean junction table: `id`, `condo_id` (FK → `condo`), `car_park_id` (no FK constraint in DB). No audit columns.

### Layer 1 — Domain model
File: `internal/vehicle/condo_car_park.go`

```go
type CondoCarPark struct {
    ID          int64  `json:"id"`
    CondoID     int64  `json:"condo_id"`
    CondoName   string `json:"condo_name"`   // display-only, from JOIN
    CarParkID   int64  `json:"car_park_id"`
    CarParkName string `json:"car_park_name"` // display-only, from JOIN
}
```

No audit fields (`created_by`, `created_at`, etc.) because the table has none.

### Layer 2 — Repository
File: `internal/vehicle/condo_car_park_repository.go`

- Interface: `CondoCarParkRepository` with `GetAll`, `GetPaged`, `Count`, `GetByID`, `Create`, `Update`, `Delete`.
- Implementation: `mysqlCondoCarParkRepository`.
- Constructor: `NewCondoCarParkRepository(db *sql.DB) CondoCarParkRepository`.
- SELECT JOINs `condo` and `car_park` to populate display names.
- Sortable columns: `id`, `condo_id`, `car_park_id` (allowlist-guarded).
- `Create` uses `LastInsertId` to populate `c.ID` after insert.

### Layer 3 — Service
File: `internal/vehicle/condo_car_park_service.go`

- Interface: `CondoCarParkService`.
- `CreateCondoCarPark(condoID, carParkID int64)` — no audit params (table has no audit columns).
- `UpdateCondoCarPark(id, condoID, carParkID int64)` — fetches first, returns `nil, nil` if not found.
- `ListPaged` clamps `page`/`pageSize` to ≥1, computes offset.

### Layer 4 — Remote service
File: `internal/vehicle/remote_service.go` (appended at bottom)

- `remoteCondoCarParkService` struct with `client *http.Client` and `baseURL string`.
- `NewRemoteCondoCarParkService(baseURL string) CondoCarParkService`.
- All 6 methods implemented; URL base is `/api/condo-car-parks`; JSON list key is `"condo_car_parks"`.

### Layer 5 — Web handler
File: `internal/http/handlers/web/condo_car_park.go`

- `CondoCarParkHandler{svc, condoSvc, cpSvc}` — holds `CondoService` and `CarParkService` for form dropdowns.
- `NewCondoCarParkHandler(svc, condoSvc, cpSvc)`.
- All 8 methods: `List`, `CreateForm`, `Create`, `EditForm`, `Update`, `Delete`, `View`, `DeleteConfirm`.
- Form passes `"condos"` and `"parks"` slices to template; record key is `"condoCarPark"`.
- `DeleteConfirm` uses `record.CondoName + " / " + record.CarParkName` as the display name.

### Layer 6 — API handler
File: `internal/http/handlers/api/condo_car_park.go`

- `CondoCarParkAPIHandler{svc}`.
- `ListAll` returns `{"condo_car_parks": [...]}`.
- `List` returns `{"condo_car_parks": [...], "total": n}`.
- Standard `Get`, `Create`, `Update`, `Delete`.

### Layer 7 — Routes
File: `internal/http/routes/routes.go`

- Added `condoCarParkHandler *web.CondoCarParkHandler` and `condoCarParkAPI *api.CondoCarParkAPIHandler` parameters.
- Added to the nil-guard condition.
- Web group at `/condo-car-parks` with all 8 routes.
- API routes at `/api/condo-car-parks` (including `/all`).

### Layer 8 — Main wiring
File: `cmd/web/main.go`

- Added `condoCarParkHandler` and `condoCarParkAPI` to the top-level var block.
- `vehicle-service` case: wires repo → svc → API handler only (no web handler).
- `web-view` case: declares `var condoCarParkSvc vehicle.CondoCarParkService`; remote branch calls `NewRemoteCondoCarParkService`; fallback branch creates local repo/svc; handler created with `condoSvc` and `cpSvc`.
- `monolith` case: full local wiring — repo → svc → web handler + API handler.
- Both new params passed to `routes.RegisterRoutes(...)`.

---

## What's next

| Layer | File(s) | Notes |
|---|---|---|
| 9 | `templates/pages/condo_car_parks/` | `index.html`, `form.html`, `view.html` + `templates/partials/condo_car_parks/table_row.html` |
| 10 | `templates/layouts/base.html` | Add sidebar nav link |

### Key decisions for Layer 9
- **Template data keys**: list page uses `"condoCarParks"` (slice); form and view use `"condoCarPark"` (single record), `"condos"` (condo dropdown), `"parks"` (car park dropdown).
- **URL slug**: `/condo-car-parks` (already registered).
- **Row ID**: `condo-car-park-row-{{ .ID }}` (used by HTMX delete swap).
- **Sortable columns**: `id`, `condo_id`, `car_park_id`.
- **Sidebar section**: decide which section this belongs under (e.g. "Condo" or "Car Park" grouping) before Layer 10.
