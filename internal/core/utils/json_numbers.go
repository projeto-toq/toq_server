package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// ParseUintFromJSON normalizes JSON numbers that may arrive as strings or numeric literals.
func ParseUintFromJSON(raw json.RawMessage) (uint64, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return 0, fmt.Errorf("value is empty")
	}

	if trimmed[0] == '"' {
		var str string
		if err := json.Unmarshal(raw, &str); err != nil {
			return 0, err
		}
		trimmed = bytes.TrimSpace([]byte(str))
		if len(trimmed) == 0 {
			return 0, fmt.Errorf("value is empty")
		}
	}

	decoder := json.NewDecoder(bytes.NewReader(trimmed))
	decoder.UseNumber()

	token, err := decoder.Token()
	if err != nil {
		return 0, err
	}

	number, ok := token.(json.Number)
	if !ok {
		return 0, fmt.Errorf("value is not a number")
	}

	value, err := strconv.ParseUint(number.String(), 10, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}
