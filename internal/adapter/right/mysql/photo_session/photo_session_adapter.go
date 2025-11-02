package mysqlphotosessionadapter

import (
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
)

// PhotoSessionAdapter provides DB access to photographer agenda entries and bookings.
type PhotoSessionAdapter struct {
	mysqladapter.InstrumentedAdapter
}

// NewPhotoSessionAdapter builds a new adapter instance.
func NewPhotoSessionAdapter(db *mysqladapter.Database, metrics metricsport.MetricsPortInterface) *PhotoSessionAdapter {
	return &PhotoSessionAdapter{InstrumentedAdapter: mysqladapter.NewInstrumentedAdapter(db, metrics)}
}
