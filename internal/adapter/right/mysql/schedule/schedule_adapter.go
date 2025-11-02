package mysqlscheduleadapter

import (
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
)

// ScheduleAdapter provides DB access to listing agendas.
type ScheduleAdapter struct {
	mysqladapter.InstrumentedAdapter
}

// NewScheduleAdapter builds a new adapter.
func NewScheduleAdapter(db *mysqladapter.Database, metrics metricsport.MetricsPortInterface) *ScheduleAdapter {
	return &ScheduleAdapter{InstrumentedAdapter: mysqladapter.NewInstrumentedAdapter(db, metrics)}
}
