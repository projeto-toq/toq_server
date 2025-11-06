package mysqlvisitadapter

import (
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
)

// VisitAdapter provides persistence access for listing visits
//
// This adapter implements the VisitRepositoryInterface port, handling all
// database operations for the listing_visits table. It uses the InstrumentedAdapter
// to automatically generate metrics, logging, and distributed tracing for all queries.
//
// Responsibilities:
//   - CRUD operations for visit scheduling (insert, get, update, list)
//   - Dynamic filtering and pagination for visit listings
//   - Entity ↔ Domain model conversions (via converters package)
//   - Query instrumentation (metrics + tracing via InstrumentedAdapter)
//
// Port Implementation:
//   - Interface: visitrepository.VisitRepositoryInterface
//   - Location: internal/core/port/right/repository/visit_repository/
//
// Database Table:
//   - Name: listing_visits
//   - Engine: InnoDB (supports transactions)
//   - Charset: utf8mb4_unicode_ci
//   - Key columns: id (PK), listing_id (FK), owner_id, realtor_id
//
// Transaction Support:
//   - All methods accept *sql.Tx parameter for transaction participation
//   - Operations MUST run within transactions for data consistency
//   - Transactions managed by service layer (global_services/transactions)
//
// Observability:
//   - Metrics: Query duration, error rate, rows affected (via InstrumentedAdapter)
//   - Tracing: Distributed tracing spans for each operation (OpenTelemetry)
//   - Logging: Contextual logging with request_id/trace_id (slog)
//
// Error Handling:
//   - Returns pure error types (no HTTP coupling)
//   - sql.ErrNoRows returned when record not found (service maps to 404)
//   - Database errors wrapped with context: fmt.Errorf("operation: %w", err)
//
// Usage:
//
//	adapter := NewVisitAdapter(database, metricsAdapter)
//	visit, err := adapter.GetVisitByID(ctx, tx, visitID)
type VisitAdapter struct {
	mysqladapter.InstrumentedAdapter
}

// NewVisitAdapter creates a new VisitAdapter instance with instrumentation
//
// This constructor initializes the adapter with database connection and metrics
// provider, enabling automatic query instrumentation for observability.
//
// Parameters:
//   - db: Database connection pool (*mysqladapter.Database)
//   - metrics: Metrics interface for recording query performance (MetricsPortInterface)
//
// Returns:
//   - *VisitAdapter: Fully initialized adapter ready for repository operations
//
// The InstrumentedAdapter provides:
//   - ExecContext: Instrumented version of db.ExecContext/tx.ExecContext
//   - QueryContext: Instrumented version of db.QueryContext/tx.QueryContext
//   - QueryRowContext: Instrumented version of db.QueryRowContext/tx.QueryRowContext
//   - Automatic metrics collection (query duration, errors, rows affected)
//   - Contextual logging with trace_id propagation
//
// Lifecycle:
//   - Created during bootstrap phase 04 (dependency injection)
//   - Registered in factory.CreateRepositoryAdapters()
//   - Shared across all service instances (singleton pattern)
func NewVisitAdapter(db *mysqladapter.Database, metrics metricsport.MetricsPortInterface) *VisitAdapter {
	return &VisitAdapter{InstrumentedAdapter: mysqladapter.NewInstrumentedAdapter(db, metrics)}
}
