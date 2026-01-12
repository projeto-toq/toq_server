package mysqllistingfavoriteadapter

import (
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
)

// ListingFavoriteAdapter implements FavoriteRepoPortInterface using MySQL with instrumentation.
type ListingFavoriteAdapter struct {
	mysqladapter.InstrumentedAdapter
}

// NewListingFavoriteAdapter builds a new adapter wired with metrics/tracing.
func NewListingFavoriteAdapter(db *mysqladapter.Database, metrics metricsport.MetricsPortInterface) *ListingFavoriteAdapter {
	return &ListingFavoriteAdapter{InstrumentedAdapter: mysqladapter.NewInstrumentedAdapter(db, metrics)}
}
