# Handover — Vehicle Global Setting

## Feature
CRUD management for the `global_setting` table in the `vehicles` database.

## Naming conventions
| Context | Name |
|---|---|
| DB table | `global_setting` |
| Go types / files | `VehicleGlobalSetting` / `vehicle_global_setting_*` |
| UI label (visible to user) | **Vehicle Service Settings** |
| URL slug | `/vehicle-global-settings` |
| API base path | `/api/vehicle-global-settings` |
| JSON list key | `settings` |

## Session 1 — completed (Layers 1–3)

### Files created
- `internal/vehicle/vehicle_global_setting.go` — domain struct
- `internal/vehicle/vehicle_global_setting_repository.go` — repository interface + MySQL impl
- `internal/vehicle/vehicle_global_setting_service.go` — service interface + impl

### Non-obvious decisions
- **`key` is a MySQL reserved word** — every SQL statement wraps it in backticks (`` `key` ``). The Go field is `Key string`.
- `remark`, `country_code`, `updated_by`, `updated_at` are nullable → pointer types (`*string`, `*time.Time`).
- `old_id` does not exist in this table — no exclusion needed.
- Sortable columns allowlist: `id`, `key`, `value`, `country_code`, `created_at`, `updated_at`.
- Schema DDL saved at `schema/vehicles/global_setting.sql`.

## Session 2 — completed (Layers 4–6)

### Files created / modified
- `internal/vehicle/remote_service.go` — appended `remoteVehicleGlobalSettingService` at the bottom
- `internal/http/handlers/web/vehicle_global_setting.go` — web handler (new file)
- `internal/http/handlers/api/vehicle_global_setting.go` — API handler (new file)

### Non-obvious decisions
- **Delete confirm name** uses `setting.Key` (the natural human identifier for a setting row).
- **HTMX row ID** is `vehicle-global-setting-row-{id}` — used by the delete swap.
- **`optionalPostForm` helper** defined in the web handler file (not in constants.go) since it's only used by this handler.
- Web handler follows AGENTS.md strictly: every 5xx path calls `slog.Error` with a safe message (no `err.Error()` exposed to the client).
- Optional fields (`remark`, `country_code`) parsed from form as `*string` — nil when the field is empty string.

## Session 3 — completed (Layers 7–8)

### Files modified
- `internal/http/routes/routes.go` — added `vgsHandler`/`vgsAPI` params, nil-guard, `/vehicle-global-settings` web routes, `/api/vehicle-global-settings` API routes
- `cmd/web/main.go` — declared `vgsHandler`/`vgsAPI` vars; wired local repo→service→handlers in `vehicle-service` and `monolith` cases; wired remote service in `web-view` (both remote and fallback-local branches); passed both to `RegisterRoutes`

### Non-obvious decisions
- `vehicle-service` case wires API handler only (no web handler) — consistent with how all other vehicle entities behave in that role.
- `web-view` case wires only the web handler (no API) — again consistent with existing pattern.
- `monolith` case wires both web handler and API handler.

## Session 4 — completed (Layers 9–10)

### Files created
- `templates/pages/vehicle_global_settings/index.html` — paginated table: ID, Key, Value, Country Code, Updated At; timezone selector; sort on all except Actions
- `templates/partials/vehicle_global_settings/table_row.html` — `<tr id="vehicle-global-setting-row-{{ .Setting.ID }}">` with view/edit/delete action buttons
- `templates/pages/vehicle_global_settings/form.html` — create/edit form: Key (required, autofocus), Value (required), Remark (optional), Country Code (optional)
- `templates/pages/vehicle_global_settings/view.html` — read-only mirror of form + form-meta audit block; Back to List footer

### Non-obvious decisions
- Nullable pointer fields (`Remark`, `CountryCode`) rendered with `{{ if .setting.Remark }}...{{ end }}` — Go templates auto-dereference `*string`.
- Index shows 5 data columns (ID, Key, Value, Country Code, Updated At); Updated By omitted from the table row to keep the row compact (all detail is in the view page).
- Timezone selector and auto-detect script included (consistent with vehicle_colors pattern).

## Session 5 — completed (Layer 11)

### Files modified
- `templates/layouts/base.html` — added "Service Settings" link (`/vehicle-global-settings`, `fa-sliders` icon) at the bottom of the Vehicle Management section

## All sessions complete — feature fully built
