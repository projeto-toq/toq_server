package mysqlcomplexadapter

import (
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
)

type ComplexAdapter struct {
	mysqladapter.InstrumentedAdapter
}

func NewComplexAdapter(db *mysqladapter.Database, metrics metricsport.MetricsPortInterface) *ComplexAdapter {
	return &ComplexAdapter{InstrumentedAdapter: mysqladapter.NewInstrumentedAdapter(db, metrics)}
}
