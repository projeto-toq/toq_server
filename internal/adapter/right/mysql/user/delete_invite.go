package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteInviteByID removes an agency invitation by its unique identifier
//
// This function deletes a specific agency_invites record. This is typically
// called when an invite is accepted by the realtor, rejected, or manually
// revoked by the agency.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (REQUIRED - deletion must be transactional)
//   - id: Primary key of the agency_invites record to delete
//
// Returns:
//   - deleted: Number of rows deleted (0 or 1)
//   - error: sql.ErrNoRows if invite not found, database errors otherwise
//
// Business Rules:
//   - Returns sql.ErrNoRows if invite ID does not exist (standard not-found behavior)
//   - Service layer maps sql.ErrNoRows to domain error (404 Not Found)
//   - Does NOT check if invite was already used (service layer responsibility)
//
// Usage Example:
//
//	deleted, err := adapter.DeleteInviteByID(ctx, tx, inviteID)
//	if err == sql.ErrNoRows {
//	    return derrors.NotFound("Invite not found")
//	}
func (ua *UserAdapter) DeleteInviteByID(ctx context.Context, tx *sql.Tx, id int64) (deleted int64, err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Delete invite record by primary key
	query := `DELETE FROM agency_invites WHERE id = ?;`

	// Execute deletion using instrumented adapter
	result, execErr := ua.ExecContext(ctx, tx, "delete", query, id)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.delete_invite.exec_error", "error", execErr)
		return 0, fmt.Errorf("delete invite by id: %w", execErr)
	}

	// Check if invite was found and deleted
	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.delete_invite.rows_affected_error", "error", rowsErr)
		return 0, fmt.Errorf("delete invite rows affected: %w", rowsErr)
	}

	// Return sql.ErrNoRows if invite not found (standard repository pattern)
	// ✅ CHANGED: Previously returned errors.New(), now returns sql.ErrNoRows
	if rowsAffected == 0 {
		return 0, sql.ErrNoRows
	}

	return rowsAffected, nil
}
