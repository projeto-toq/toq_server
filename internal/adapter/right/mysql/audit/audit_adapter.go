package mysqlauditadapter

import (
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
)

// AuditAdapter persists audit events into MySQL using the instrumented adapter.
type AuditAdapter struct {
	mysqladapter.InstrumentedAdapter
}

// NewAuditAdapter builds a new AuditAdapter instance.
func NewAuditAdapter(db *mysqladapter.Database, metrics metricsport.MetricsPortInterface) *AuditAdapter {
	return &AuditAdapter{InstrumentedAdapter: mysqladapter.NewInstrumentedAdapter(db, metrics)}
}
