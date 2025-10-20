package mysqlcomplexadapter

import "strings"

func nullableStringValue(value string) any {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return trimmed
}

func nullableIntValue(value *int) any {
	if value == nil {
		return nil
	}
	return *value
}
