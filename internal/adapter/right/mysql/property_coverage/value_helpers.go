package mysqlpropertycoverageadapter

import "strings"

// // nullableString returns nil for blank strings so nullable DB columns can be set correctly.
// func nullableString(value string) any {
// 	trimmed := strings.TrimSpace(value)
// 	if trimmed == "" {
// 		return nil
// 	}
// 	return trimmed
// }

// likePattern builds a basic LIKE pattern with trimmed value wrapped in % wildcards.
func likePattern(value string) string {
	return "%" + strings.TrimSpace(value) + "%"
}

// // pointerIntOrZero unwraps an *int, returning zero when the pointer is nil.
// func pointerIntOrZero(value *int) int {
// 	if value == nil {
// 		return 0
// 	}
// 	return *value
// }
