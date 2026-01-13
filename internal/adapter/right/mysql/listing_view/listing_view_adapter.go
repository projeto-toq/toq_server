package mysqllistingviewadapter

import (
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
)

// ListingViewAdapter provides MySQL persistence for listing view counters.
// It embeds InstrumentedAdapter to leverage unified tracing, metrics, and logging.
type ListingViewAdapter struct {
	mysqladapter.InstrumentedAdapter
}

// NewListingViewAdapter builds a new ListingViewAdapter with instrumentation enabled.
func NewListingViewAdapter(db *mysqladapter.Database, metrics metricsport.MetricsPortInterface) *ListingViewAdapter {
	return &ListingViewAdapter{InstrumentedAdapter: mysqladapter.NewInstrumentedAdapter(db, metrics)}
}
