package mysqllistingadapter

import (
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
)

type ListingAdapter struct {
	mysqladapter.InstrumentedAdapter
}

func NewListingAdapter(db *mysqladapter.Database, metrics metricsport.MetricsPortInterface) *ListingAdapter {
	return &ListingAdapter{InstrumentedAdapter: mysqladapter.NewInstrumentedAdapter(db, metrics)}
}
