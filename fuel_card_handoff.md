# fuel_card Feature — Handoff Context

## What this document is

This file is a handoff context for another AI session continuing the `fuel_card` CRUD feature. Layers 1–8 are complete and the build is clean. This document covers what was done, the key decisions made, and exactly what still needs to be built (layers 9–10).

---

## Database table

Database: `vehicles`  
Table: `fuel_card`

```sql
CREATE TABLE `fuel_card` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `card_no` varchar(255) NOT NULL,
  `fuel_company_id` bigint NOT NULL,   -- FK → fuel_company.id (ON DELETE CASCADE)
  `pin_number` varchar(255) NOT NULL,
  `vehicle_id` bigint DEFAULT NULL,    -- nullable FK → vehicle.id
  `status` tinyint(1) NOT NULL,        -- bool
  `created_by` varchar(50) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_by` varchar(50) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `source_db` varchar(10) DEFAULT NULL, -- ignored in app code (migration artifact)
  `old_id` varchar(255) DEFAULT NULL,   -- ignored in app code (legacy artifact)
  PRIMARY KEY (`id`)
)
```

---

## Design decisions already made

| Decision | Choice |
|---|---|
| `source_db` | Ignored in application code (same treatment as `old_id`) |
| `pin_number` | Included as a plain `string` in the domain struct |
| `vehicle_id` FK join display | Join `vehicle` on `vehicle_id = vehicle.id`, show `plate_number` (nullable → `*string`) |
| `fuel_company_id` FK join display | Join `fuel_company`, show `name` as `FuelCompanyName string` |

---

## Files already created

| File | Layer | Status |
|---|---|---|
| `schema/vehicles/fuel_card.sql` | DDL | Done |
| `internal/vehicle/fuel_card.go` | Layer 1 — Domain model | Done |
| `internal/vehicle/fuel_card_repository.go` | Layer 2 — Repository | Done |
| `internal/vehicle/fuel_card_service.go` | Layer 3 — Service | Done |
| `internal/vehicle/remote_service.go` (appended) | Layer 4 — Remote service | Done |
| `internal/http/handlers/web/fuel_card.go` | Layer 5 — Web handler | Done |
| `internal/http/handlers/api/fuel_card.go` | Layer 6 — API handler | Done |
| `internal/http/routes/routes.go` (updated) | Layer 7 — Routes | Done |
| `cmd/web/main.go` (updated) | Layer 8 — Main wiring | Done |

---

## Domain model (`internal/vehicle/fuel_card.go`)

```go
type FuelCard struct {
    ID              int64      `json:"id"`
    CardNo          string     `json:"card_no"`
    FuelCompanyID   int64      `json:"fuel_company_id"`
    FuelCompanyName string     `json:"fuel_company_name"`  // display join
    PinNumber       string     `json:"pin_number"`
    VehicleID       *int64     `json:"vehicle_id"`         // nullable
    PlateNumber     *string    `json:"plate_number"`       // display join, nullable
    Status          bool       `json:"status"`
    CreatedBy       string     `json:"created_by"`
    CreatedAt       time.Time  `json:"created_at"`
    UpdatedBy       *string    `json:"updated_by"`
    UpdatedAt       *time.Time `json:"updated_at"`
}
```

## Service interface (`internal/vehicle/fuel_card_service.go`)

```go
type FuelCardService interface {
    ListAll() ([]FuelCard, error)
    ListPaged(page, pageSize int, sortBy, sortOrder string) ([]FuelCard, int, error)
    FindByID(id int64) (*FuelCard, error)
    CreateFuelCard(cardNo string, fuelCompanyID int64, pinNumber string, vehicleID *int64, status bool, user string) (*FuelCard, error)
    UpdateFuelCard(id int64, cardNo string, fuelCompanyID int64, pinNumber string, vehicleID *int64, status bool, user string) (*FuelCard, error)
    DeleteFuelCard(id int64) error
}
```

Constructor: `NewFuelCardService(repo FuelCardRepository) FuelCardService`

## Web handler (`internal/http/handlers/web/fuel_card.go`)

```go
type FuelCardHandler struct {
    svc   vehicle.FuelCardService
    fcSvc vehicle.FuelCompanyService   // for fuel company dropdown
    vSvc  vehicle.VehicleService       // for vehicle dropdown
}

func NewFuelCardHandler(svc vehicle.FuelCardService, fcSvc vehicle.FuelCompanyService, vSvc vehicle.VehicleService) *FuelCardHandler
```

- `CreateForm`/`EditForm` pass `"companies"` (`[]FuelCompany`) and `"vehicles"` (`[]Vehicle`) to the template.
- Nullable `vehicle_id`: empty string or `"0"` from the form is treated as `nil`.
- Template data keys: `"card"` (single record), `"cards"` (list), `"companies"`, `"vehicles"`.

## API handler (`internal/http/handlers/api/fuel_card.go`)

```go
type FuelCardAPIHandler struct{ svc vehicle.FuelCardService }
func NewFuelCardAPIHandler(svc vehicle.FuelCardService) *FuelCardAPIHandler
```

`ListAll` returns `{"cards": [...]}`. `List` returns `{"cards": [...], "total": N}`.

## Routes registered

Web routes (protected, under `/fuel-cards`):
```
GET    /fuel-cards
GET    /fuel-cards/new
POST   /fuel-cards
GET    /fuel-cards/:id/view
GET    /fuel-cards/:id/edit
POST   /fuel-cards/:id
DELETE /fuel-cards/:id
GET    /fuel-cards/:id/delete
```

API routes (under `/api`):
```
GET    /api/fuel-cards/all
GET    /api/fuel-cards
GET    /api/fuel-cards/:id
POST   /api/fuel-cards
PUT    /api/fuel-cards/:id
DELETE /api/fuel-cards/:id
```

## Repository interface (`internal/vehicle/fuel_card_repository.go`)

Constructor: `NewFuelCardRepository(db *sql.DB) FuelCardRepository`

Sortable columns: `id`, `card_no`, `status`, `updated_by`, `updated_at` (prefixed `fc.` in ORDER BY).

---

## What still needs to be built (layers 9–10)

Follow the scaffold in `AGENTS.md` exactly. Use `FuelCompanyHandler` / `FuelCompanyAPIHandler` in the same files as the canonical reference for a feature with audit fields + FK dropdown — those files are at:
- `internal/http/handlers/web/fuel_company.go`
- `internal/http/handlers/api/fuel_company.go`

Canonical template reference: `templates/pages/fuel_companies/` for audit-field + status + FK dropdown pattern.  
Car park lot form reference: `templates/pages/car_park_lots/form.html` for the FK `<select>` dropdown pattern.

### Layer 9 — Templates

| File | Purpose |
|---|---|
| `templates/pages/fuel_cards/index.html` | Paginated table with sort, rows-per-page, timezone selector, pagination nav |
| `templates/pages/fuel_cards/form.html` | Create/edit form; dropdowns for fuel company (required) and vehicle (optional) |
| `templates/pages/fuel_cards/view.html` | Read-only mirror of the form |
| `templates/partials/fuel_cards/table_row.html` | `<tr id="fuel-card-row-{{ .ID }}">` with view/edit/delete action buttons |

Template data available to each template:

**index.html** — `.cards` `[]FuelCard`, `.timezone`, `.timezones`, `.currentPage`, `.pageSize`, `.pageSizeOptions`, `.totalPages`, `.total`, `.start`, `.end`, `.sortBy`, `.sortOrder`, `.pageWindow`

**form.html** — `.title`, `.action`, `.card` (nil when creating), `.companies` `[]FuelCompany`, `.vehicles` `[]Vehicle`

**view.html** — `.title`, `.card` `*FuelCard`, `.tz`

**table_row.html** — called as `{{ template "fuel_cards/table_row.html" (dict "Card" . "Timezone" $tz) }}`

Template notes:
- Status field: render as a pill/badge (`pill pill-on` / `pill pill-off`), consistent with `fuel_companies`.
- Index table columns: ID, Card No, Fuel Company, Vehicle, Status, Updated By, Updated At, Actions.
- Sortable columns (match repo allowlist): `id`, `card_no`, `status`, `updated_by`, `updated_at`.
- `vehicle_id` in the form: a `<select>` with an empty `<option value="0">— Unassigned —</option>` at the top (FK is nullable). Show `plate_number` as the display label; fall back to the vehicle ID if `plate_number` is nil. Pre-select if `.card.VehicleID` is non-nil.
- `fuel_company_id` in the form: required `<select>`; pre-select if `.card.FuelCompanyID` matches.
- `pin_number`: standard text input. No masking required.
- Autofocus on the first input (`card_no`).
- Pagination footer uses `.start`/`.end`/`.total` and `.pageWindow` (render `-1` entries as `<span class="pag-ellipsis">…</span>`).

### Layer 10 — Sidebar nav

File: `templates/layouts/base.html`

Add a link **after** the Fuel Companies entry (line ~79) under the Vehicle Management `<div class="sb-section">`:

```html
<a href="/fuel-cards" class="sb-item" data-href="/fuel-cards">
    <span class="sb-ico"><i class="fa-solid fa-credit-card fa-sm"></i></span>
    Fuel Cards
</a>
```

---

## Project conventions reminder

- Package: `vehicle` (all domain files share `package vehicle`)
- All timestamps UTC in DB; display in user-selected timezone in UI
- `slog.Error(...)` before every 5xx return in handlers; never pass `err.Error()` to client
- `paginationWindow(page, totalPages int) []int` helper exists in `internal/http/handlers/web/constants.go` — use it for the index page pagination
- Search: button-triggered only, disabled at exactly 1 character, enabled at 0 or 2+
- AGENTS.md is the authoritative guide — read it for any convention not covered here
