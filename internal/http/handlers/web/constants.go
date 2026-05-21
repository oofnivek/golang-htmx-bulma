package web

// Timezones provides a shared list of supported timezones for UI rendering.
// The slice is exported so that other handlers can reference it.
var Timezones = []string{"UTC", "America/New_York", "America/Los_Angeles", "Europe/London", "Europe/Paris", "Asia/Tokyo", "Asia/Shanghai", "Asia/Singapore", "Australia/Sydney"}
