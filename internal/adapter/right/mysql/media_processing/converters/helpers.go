package mediaprocessingconverters

import (
	"database/sql"
	"encoding/json"
	"time"
)

func encodeStringMap(data map[string]string) sql.NullString {
	if len(data) == 0 {
		return sql.NullString{}
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return sql.NullString{}
	}
	return sql.NullString{String: string(payload), Valid: true}
}

// EncodeStatusDetails expõe o helper para outros packages do adapter.
func EncodeStatusDetails(data map[string]string) sql.NullString {
	return encodeStringMap(data)
}

func decodeStringMap(raw sql.NullString) map[string]string {
	if !raw.Valid || raw.String == "" {
		return map[string]string{}
	}
	var result map[string]string
	if err := json.Unmarshal([]byte(raw.String), &result); err != nil {
		return map[string]string{}
	}
	return result
}

// DecodeStatusDetails expõe a conversão reversa.
func DecodeStatusDetails(raw sql.NullString) map[string]string {
	return decodeStringMap(raw)
}

func nullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: value, Valid: true}
}

func nullTimeFromPtr(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func timePtrFromNull(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	result := value.Time
	return &result
}
