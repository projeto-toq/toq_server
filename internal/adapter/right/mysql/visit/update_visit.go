package mysqlvisitadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/visit/converters"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateVisit updates an existing visit record in the listing_visits table.
// Returns sql.ErrNoRows if the visit ID does not exist.
//
// This function performs a full row update (all fields except id) and must run
// within a transaction to ensure atomicity with related operations like audit logging.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (must not be nil for consistency)
//   - visit: VisitInterface with ID populated and all fields to be updated
//
// Returns:
//   - error: sql.ErrNoRows if visit not found, or database/constraint errors
//
// Update Scope:
//   - Updates all fields: listing_id, owner_id, realtor_id, timestamps, status, notes
//   - Does NOT update: id (immutable primary key), created_by (audit immutability)
//   - Sets updated_by to current user (must be set in service layer before calling)
//
// Database Constraints Enforced:
//   - Foreign key constraints on listing_id (must reference existing listing)
//   - Status must be valid ENUM value
//   - Referential integrity maintained for owner_id and realtor_id
//
// Error Scenarios:
//   - sql.ErrNoRows: Visit ID does not exist (RowsAffected = 0)
//   - Foreign key violation: Invalid listing_id after update
//   - Enum violation: Invalid status transition
//   - Transaction conflict: Concurrent update detected
func (a *VisitAdapter) UpdateVisit(ctx context.Context, tx *sql.Tx, visit listingmodel.VisitInterface) error {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	// Ensure logger propagation with request_id and trace_id
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Convert domain model to database entity
	entity := converters.ToVisitEntity(visit)

	// UPDATE query with all mutable fields
	// Note: WHERE id = ? ensures we only update the target visit
	query := `UPDATE listing_visits SET listing_id = ?, owner_id = ?, realtor_id = ?, scheduled_start = ?, scheduled_end = ?, status = ?, cancel_reason = ?, notes = ?, updated_by = ? WHERE id = ?`

	// Execute update via instrumented adapter
	result, err := a.ExecContext(ctx, tx, "update_visit", query,
		entity.ListingID,
		entity.OwnerID,
		entity.RealtorID,
		entity.ScheduledStart,
		entity.ScheduledEnd,
		entity.Status,
		entity.CancelReason,
		entity.Notes,
		entity.UpdatedBy,
		entity.ID,
	)
	if err != nil {
		// Mark span as error for distributed tracing
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.visit.update.exec_error", "visit_id", entity.ID, "err", err)
		return fmt.Errorf("update visit: %w", err)
	}

	// Check if any rows were affected (does visit exist?)
	affected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.visit.update.rows_error", "visit_id", entity.ID, "err", err)
		return fmt.Errorf("visit rows affected: %w", err)
	}

	// Return sql.ErrNoRows if visit not found (service layer maps to 404)
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
