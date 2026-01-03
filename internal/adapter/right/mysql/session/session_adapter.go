package sessionmysqladapter

import (
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
	sessionrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/session_repository"
)

// SessionAdapter implements SessionRepoPortInterface using MySQL with full observability
//
// This adapter mirrors the same documentation depth used in the user adapter. It owns
// all database access for sessions while delegating conversions to dedicated converters
// and instrumentation to InstrumentedAdapter.
//
// Responsibilities:
//   - Execute SQL operations for session lifecycle (create, read, rotate, revoke, delete)
//   - Convert between database entities and domain models (via converters package)
//   - Preserve hexagonal separation: no business logic, no HTTP concerns
//   - Keep one public method per file per project guide for readability
//
// Architecture:
//   - Implements: internal/core/port/right/repository/session_repository.SessionRepoPortInterface
//   - Uses: InstrumentedAdapter for ExecContext/QueryContext/QueryRowContext
//   - Delegates: DB ↔ domain translations to internal/adapter/right/mysql/session/converters
//   - Shared row mapping: mapSessionFromScanner centralizes scan logic
//
// Observability:
//   - Tracing: Each public method starts a span via utils.GenerateTracer
//   - Logging: Contextual logger (utils.LoggerFromContext) propagates request_id/trace_id
//   - Metrics: InstrumentedAdapter auto-records duration/error metrics per operation name
//
// Transaction Handling:
//   - All methods accept *sql.Tx (nil allowed for standalone operations)
//   - Transaction lifecycle is owned by the service layer; adapter always reuses provided tx
//
// Error Handling:
//   - Returns sql.ErrNoRows for not-found lookups; callers map to HTTP/domain errors
//   - Wraps infrastructure errors with context (fmt.Errorf("...: %w", err))
//   - Marks spans and logs on infrastructure failures
//
// Usage Example:
//
//	adapter := NewSessionAdapter(db, metrics)
//	session := sessionmodel.NewSession()
//	session.SetUserID(123)
//	session.SetRefreshHash("abc...")
//	_ = adapter.CreateSession(ctx, tx, session)
//
// Ensure implementation satisfies port interface
var _ sessionrepository.SessionRepoPortInterface = (*SessionAdapter)(nil)

type SessionAdapter struct {
	mysqladapter.InstrumentedAdapter
}

// NewSessionAdapter builds a SessionAdapter with instrumentation (metrics + tracing + contextual logging).
//
// Parameters:
//   - db: MySQL database wrapper created in bootstrap
//   - metrics: Prometheus metrics port for instrumentation
//
// Returns:
//   - SessionRepoPortInterface: Concrete adapter implementation
//
// Lifecycle:
//   - Created once during dependency injection (Phase 04)
//   - Safe for concurrent use; holds no per-request state
//   - Database pool lifecycle is managed externally by bootstrap
func NewSessionAdapter(db *mysqladapter.Database, metrics metricsport.MetricsPortInterface) sessionrepository.SessionRepoPortInterface {
	return &SessionAdapter{
		InstrumentedAdapter: mysqladapter.NewInstrumentedAdapter(db, metrics),
	}
}
