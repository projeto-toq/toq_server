package mysqlownermetricsadapter

import (
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
)

// OwnerMetricsAdapter implements the owner metrics repository backed by MySQL.
type OwnerMetricsAdapter struct {
	mysqladapter.InstrumentedAdapter
}

// NewOwnerMetricsAdapter wires the instrumented adapter for SLA aggregation queries.
func NewOwnerMetricsAdapter(
	db *mysqladapter.Database,
	metrics metricsport.MetricsPortInterface,
) *OwnerMetricsAdapter {
	return &OwnerMetricsAdapter{
		InstrumentedAdapter: mysqladapter.NewInstrumentedAdapter(db, metrics),
	}
}
