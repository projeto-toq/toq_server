package mysqlvisitadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/visit/converters"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetVisitByID retrieves a visit by its unique ID from the listing_visits table.
// Returns sql.ErrNoRows if no visit is found with the given ID.
//
// This function fetches a single visit record without any soft delete filtering
// (table does not have a 'deleted' column). All visits are retrievable regardless
// of their status (PENDING, APPROVED, REJECTED, CANCELLED, COMPLETED, NO_SHOW).
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging. Must contain request metadata.
//   - tx: Database transaction (can be nil for standalone queries outside transactions)
//   - id: Visit's unique identifier (INT UNSIGNED PRIMARY KEY)
//
// Returns:
//   - visit: VisitInterface containing all visit data including scheduling and status
//   - error: sql.ErrNoRows if not found, or other database/scan errors
//
// Query Details:
//   - Table: listing_visits
//   - Conditions: id = ? (no status or deletion filters)
//   - Performance: Uses primary key lookup (O(1) with index)
//
// Error Scenarios:
//   - sql.ErrNoRows: Visit ID does not exist
//   - Scan errors: Database type mismatch or corruption
//   - Context cancellation: Query interrupted
func (a *VisitAdapter) GetVisitByID(ctx context.Context, tx *sql.Tx, id int64) (listingmodel.VisitInterface, error) {
	// Initialize tracing for observability (distributed tracing + metrics)
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	// Attach logger to context to ensure request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query selects all columns from listing_visits for the specified ID
	// No status filter: all visit states are retrievable (business rule)
	query := `SELECT id, listing_identity_id, listing_version, user_id, owner_user_id, scheduled_start, scheduled_end, duration_minutes, status, type, source, realtor_notes, owner_notes, rejection_reason, cancel_reason, first_owner_action_at FROM listing_visits WHERE id = ?`

	// Execute query using instrumented adapter (auto-generates metrics + tracing)
	row := a.QueryRowContext(ctx, tx, "get_visit_by_id", query, id)

	// Convert database row to strongly-typed entity
	visitEntity, err := scanVisitEntity(row)
	if err != nil {
		// Handle no results: return standard sql.ErrNoRows for service layer handling
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}

		// Mark span as error for distributed tracing analysis
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.visit.get.scan_error", "visit_id", id, "err", err)
		return nil, fmt.Errorf("scan visit: %w", err)
	}

	// Convert database entity to domain model (separation of concerns)
	return converters.ToVisitModel(visitEntity), nil
}
