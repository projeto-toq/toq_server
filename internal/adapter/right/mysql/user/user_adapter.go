package mysqluseradapter

import (
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
)

// UserAdapter implements UserRepoPortInterface using MySQL as the persistence layer
//
// This adapter provides all user-related database operations following the Repository pattern.
// It uses InstrumentedAdapter to ensure all database operations are automatically traced,
// logged, and metered for observability.
//
// Responsibilities:
//   - Execute SQL queries for user CRUD operations
//   - Convert between database entities and domain models (via converters)
//   - Return pure errors (sql.ErrNoRows, fmt.Errorf) for service layer mapping
//   - Maintain separation of concerns (no business logic, no HTTP dependencies)
//
// Architecture:
//   - Implements: internal/core/port/right/repository/user_repository.UserRepoPortInterface
//   - Uses: InstrumentedAdapter for all DB operations (ExecContext, QueryContext, QueryRowContext)
//   - Delegates: Type conversions to internal/adapter/right/mysql/user/converters
//
// Observability:
//   - Tracing: All public methods initialize spans via utils.GenerateTracer
//   - Logging: Uses contextual logger (utils.LoggerFromContext) with request_id/trace_id
//   - Metrics: Automatic query duration, error rates via InstrumentedAdapter
//
// Database Schema:
//   - Primary table: users (see entities/user_entity.go for full schema)
//   - Related tables: user_roles, user_validations, temp_user_validations, realtors_agency
//
// Transaction Handling:
//   - Most methods accept *sql.Tx parameter (can be nil for standalone queries)
//   - Transaction lifecycle managed by service layer
//   - Always use provided tx when not nil to maintain ACID properties
//
// Error Handling:
//   - Returns sql.ErrNoRows for not found scenarios (service maps to 404)
//   - Returns raw database errors (service wraps as infrastructure errors)
//   - Logs errors with slog.Error and marks spans with utils.SetSpanError
//
// Usage Example:
//
//	adapter := NewUserAdapter(database, metricsPort)
//	user, err := adapter.GetUserByID(ctx, tx, userID)
//	if err == sql.ErrNoRows {
//	    // Handle not found
//	} else if err != nil {
//	    // Handle infrastructure error
//	}
type UserAdapter struct {
	mysqladapter.InstrumentedAdapter
}

// NewUserAdapter creates a new UserAdapter with instrumented database access and metrics
//
// The adapter is configured with InstrumentedAdapter which provides:
//   - Automatic tracing of all SQL queries
//   - Prometheus metrics for query duration and errors
//   - Contextual logging with request_id correlation
//
// Parameters:
//   - db: MySQL database connection pool (managed by config.Bootstrap)
//   - metrics: Metrics interface for Prometheus instrumentation
//
// Returns:
//   - *UserAdapter: Configured adapter ready for use
//
// Lifecycle:
//   - Created once during application bootstrap (Phase 04 - Dependency Injection)
//   - Shared across all requests (stateless, thread-safe)
//   - Database pool managed by bootstrap, not closed by adapter
func NewUserAdapter(
	db *mysqladapter.Database,
	metrics metricsport.MetricsPortInterface,
) *UserAdapter {
	return &UserAdapter{
		InstrumentedAdapter: mysqladapter.NewInstrumentedAdapter(db, metrics),
	}
}
