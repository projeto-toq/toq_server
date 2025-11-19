package mysqlpropertycoverageadapter

import "strings"

func nullableString(value string) any {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return trimmed
}

func likePattern(value string) string {
	return "%" + strings.TrimSpace(value) + "%"
}

func pointerIntOrZero(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}
