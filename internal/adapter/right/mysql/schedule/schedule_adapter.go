package mysqlscheduleadapter

import (
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
)

// ScheduleAdapter implements ScheduleRepositoryInterface using MySQL as the persistence layer.
//
// Responsibilities:
//   - Execute SQL queries for agendas/rules/entries with tracing, metrics, and structured logging (InstrumentedAdapter).
//   - Delegate all entity↔domain transformations to scheduleconverters to keep the domain isolated from persistence.
//   - Return pure infrastructure errors only (sql.ErrNoRows, fmt.Errorf); never map to HTTP or business responses.
//
// Architecture:
//   - Implements: internal/core/port/right/repository/schedule_repository.ScheduleRepositoryInterface.
//   - Uses: InstrumentedAdapter for ExecContext, QueryContext, QueryRowContext with automatic observability.
//   - Delegates: Conversion logic to internal/adapter/right/mysql/schedule/converters.
//
// Observability:
//   - Tracing: Public methods create spans via utils.GenerateTracer; errors marked with utils.SetSpanError.
//   - Logging: Contextual logger (utils.LoggerFromContext) adds request/trace correlation on errors.
//   - Metrics: InstrumentedAdapter emits duration/error metrics through metrics port.
//
// Transaction Handling:
//   - Methods accept *sql.Tx (may be nil when allowed); transactions are controlled by the service layer.
//   - Always prefer the provided tx to ensure atomicity alongside other repository calls.
//
// Error Handling:
//   - Returns sql.ErrNoRows for not-found or RowsAffected==0 when appropriate; service maps to domain errors.
//   - Bubbles driver/errors from database without wrapping to preserve diagnosis and retry behavior.
//
// Lifecycle: created once during application bootstrap (DI) and shared safely across requests.
type ScheduleAdapter struct {
	mysqladapter.InstrumentedAdapter
}

// NewScheduleAdapter wires InstrumentedAdapter with metrics for schedule operations.
// Parameters:
//   - db: database pool managed by bootstrap.
//   - metrics: metrics port for Prometheus instrumentation.
//
// Returns: configured ScheduleAdapter ready for use.
func NewScheduleAdapter(db *mysqladapter.Database, metrics metricsport.MetricsPortInterface) *ScheduleAdapter {
	return &ScheduleAdapter{InstrumentedAdapter: mysqladapter.NewInstrumentedAdapter(db, metrics)}
}
