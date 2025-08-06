package converters

import (
	"database/sql"
	"time"
)

// TimeToNullTime converte um valor time.Time para sql.NullTime
func TimeToNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{
			Time:  t,
			Valid: false,
		}
	}
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}
