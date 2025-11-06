package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteAgencyRealtorRelation removes the relationship between a realtor and an agency
//
// This function deletes the junction record in realtors_agency table, effectively
// removing the realtor from the agency's roster. The realtor user record remains
// intact (no cascade delete on users table).
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (REQUIRED - deletion must be transactional)
//   - agencyID: ID of the agency
//   - realtorID: ID of the realtor being removed
//
// Returns:
//   - deleted: Number of rows deleted (0 or 1)
//   - error: sql.ErrNoRows if relationship not found, database errors otherwise
//
// Business Rules:
//   - Returns sql.ErrNoRows if no relationship exists (idempotent operation)
//   - Service layer maps sql.ErrNoRows to domain error (404 Not Found)
//   - Does NOT delete user records (only the relationship)
//
// Database Constraints:
//   - Composite WHERE clause ensures exact match (realtor_id AND agency_id)
//
// Usage Example:
//
//	deleted, err := adapter.DeleteAgencyRealtorRelation(ctx, tx, agencyID, realtorID)
//	if err == sql.ErrNoRows {
//	    // Handle "relationship not found"
//	}
func (ua *UserAdapter) DeleteAgencyRealtorRelation(ctx context.Context, tx *sql.Tx, agencyID int64, realtorID int64) (deleted int64, err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Delete relationship record (composite WHERE for exact match)
	query := `DELETE FROM realtors_agency WHERE realtor_id = ? AND agency_id = ?;`

	// Execute deletion using instrumented adapter
	result, execErr := ua.ExecContext(ctx, tx, "delete", query, realtorID, agencyID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.delete_agency_realtor_relation.exec_error", "error", execErr)
		return 0, fmt.Errorf("delete realtor-agency relation: %w", execErr)
	}

	// Check if relationship was found and deleted
	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.delete_agency_realtor_relation.rows_affected_error", "error", rowsErr)
		return 0, fmt.Errorf("delete realtor-agency relation rows affected: %w", rowsErr)
	}

	// Return sql.ErrNoRows if relationship not found (standard repository pattern)
	if rowsAffected == 0 {
		return 0, sql.ErrNoRows
	}

	return rowsAffected, nil
}
