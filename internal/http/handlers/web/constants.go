package web

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

// Timezones provides a shared list of supported timezones for UI rendering.
// The slice is exported so that other handlers can reference it.
var Timezones = []string{"UTC", "America/New_York", "America/Los_Angeles", "Europe/London", "Europe/Paris", "Asia/Tokyo", "Asia/Shanghai", "Asia/Singapore", "Australia/Sydney"}

// isForeignKeyConstraintError reports whether err is a MySQL FK constraint
// violation (error 1451 — cannot delete a parent row referenced by a child).
func isForeignKeyConstraintError(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1451
}

// paginationWindow returns the page numbers to render in a paginator.
// -1 is a sentinel meaning "render an ellipsis here".
// It always includes page 1, the last page, and up to 2 neighbours on each
// side of the current page, with ellipsis inserted wherever there is a gap.
func paginationWindow(current, total int) []int {
	if total <= 9 {
		pages := make([]int, total)
		for i := range pages {
			pages[i] = i + 1
		}
		return pages
	}

	visible := map[int]bool{1: true, total: true}
	for d := -2; d <= 2; d++ {
		if p := current + d; p >= 1 && p <= total {
			visible[p] = true
		}
	}

	sorted := make([]int, 0, len(visible))
	for p := 1; p <= total; p++ {
		if visible[p] {
			sorted = append(sorted, p)
		}
	}

	result := make([]int, 0, len(sorted)+2)
	for i, p := range sorted {
		if i > 0 && p-sorted[i-1] > 1 {
			result = append(result, -1)
		}
		result = append(result, p)
	}
	return result
}
