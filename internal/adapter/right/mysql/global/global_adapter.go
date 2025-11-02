package mysqlglobaladapter

import (
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
)

type GlobalAdapter struct {
	mysqladapter.InstrumentedAdapter
}

func NewGlobalAdapter(db *mysqladapter.Database, metrics metricsport.MetricsPortInterface) *GlobalAdapter {
	return &GlobalAdapter{
		InstrumentedAdapter: mysqladapter.NewInstrumentedAdapter(db, metrics),
	}
}
