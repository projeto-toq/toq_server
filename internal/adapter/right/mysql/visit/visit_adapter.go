package mysqlvisitadapter

import (
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
)

// VisitAdapter provides persistence access for listing visits.
type VisitAdapter struct {
	mysqladapter.InstrumentedAdapter
}

// NewVisitAdapter creates a new adapter.
func NewVisitAdapter(db *mysqladapter.Database, metrics metricsport.MetricsPortInterface) *VisitAdapter {
	return &VisitAdapter{InstrumentedAdapter: mysqladapter.NewInstrumentedAdapter(db, metrics)}
}
