package mysqluseradapter

import (
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
)

// UserAdapter implements UserRepoPortInterface using MySQL
// This adapter provides persistence operations for user entities
// It uses InstrumentedAdapter to ensure all database operations are traced and metered
type UserAdapter struct {
	mysqladapter.InstrumentedAdapter
}

// NewUserAdapter creates a new UserAdapter with instrumented database access and metrics
// The adapter automatically generates metrics and tracing for all database operations
//
// Parameters:
//   - db: MySQL database connection pool
//   - metrics: Metrics interface for Prometheus instrumentation
//
// Returns:
//   - *UserAdapter: Configured adapter ready for use
func NewUserAdapter(
	db *mysqladapter.Database,
	metrics metricsport.MetricsPortInterface,
) *UserAdapter {
	return &UserAdapter{
		InstrumentedAdapter: mysqladapter.NewInstrumentedAdapter(db, metrics),
	}
}
