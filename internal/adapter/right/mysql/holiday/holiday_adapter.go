package mysqlholidayadapter

import mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"

// HolidayAdapter provides access to holiday calendars tables.
type HolidayAdapter struct {
	db *mysqladapter.Database
}

// NewHolidayAdapter creates a new adapter instance.
func NewHolidayAdapter(db *mysqladapter.Database) *HolidayAdapter {
	return &HolidayAdapter{db: db}
}
