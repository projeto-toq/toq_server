package mysqlpropertycoverageadapter

import (
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
	propertycoveragerepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/property_coverage_repository"
)

// PropertyCoverageAdapter implements the repository contract backed by MySQL queries.
type PropertyCoverageAdapter struct {
	mysqladapter.InstrumentedAdapter
}

// NewPropertyCoverageAdapter wires shared DB handle and metrics collector into the adapter.
func NewPropertyCoverageAdapter(db *mysqladapter.Database, metrics metricsport.MetricsPortInterface) *PropertyCoverageAdapter {
	return &PropertyCoverageAdapter{InstrumentedAdapter: mysqladapter.NewInstrumentedAdapter(db, metrics)}
}

var _ propertycoveragerepository.RepositoryInterface = (*PropertyCoverageAdapter)(nil)
