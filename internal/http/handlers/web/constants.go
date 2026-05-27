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
