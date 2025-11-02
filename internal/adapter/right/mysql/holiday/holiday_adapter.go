package mysqlholidayadapter

import (
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
)

// HolidayAdapter provides access to holiday calendars tables.
type HolidayAdapter struct {
	mysqladapter.InstrumentedAdapter
}

// NewHolidayAdapter creates a new adapter instance.
func NewHolidayAdapter(db *mysqladapter.Database, metrics metricsport.MetricsPortInterface) *HolidayAdapter {
	return &HolidayAdapter{InstrumentedAdapter: mysqladapter.NewInstrumentedAdapter(db, metrics)}
}
