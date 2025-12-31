package mysqlvisitadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/visit/converters"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// InsertVisit creates a new visit record in the listing_visits table.
// Returns the auto-generated ID of the newly created visit.
//
// This function performs a single INSERT operation and must run within a transaction
// to ensure data consistency with related operations (e.g., audit logs, notifications).
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (must not be nil for consistency)
//   - visit: VisitInterface with all required fields populated except ID
//
// Returns:
//   - id: Auto-generated INT UNSIGNED PRIMARY KEY from database
//   - error: Database constraint violations, connection errors, or transaction errors
//
// Database Constraints Enforced:
//   - listing_id must reference existing listing (FK constraint)
//   - owner_id must reference existing user (implicit FK)
//   - realtor_id must reference existing user (implicit FK)
//   - scheduled_start < scheduled_end (validated in service layer)
//   - status must be valid ENUM value (DB enforces)
//
// Side Effects:
//   - Sets the ID on the provided visit object via SetID()
//   - Updates auto_increment counter in database
//
// Error Scenarios:
//   - Foreign key violation: Invalid listing_id/owner_id/realtor_id
//   - Enum violation: Invalid status value
//   - Transaction rollback: Upstream transaction fails
func (a *VisitAdapter) InsertVisit(ctx context.Context, tx *sql.Tx, visit listingmodel.VisitInterface) (int64, error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	// Ensure logger propagation with request_id and trace_id
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Convert domain model to database entity (handles sql.Null* conversions)
	entity := converters.ToVisitEntity(visit)

	// INSERT query with all fields except id (AUTO_INCREMENT)
	// Note: created_by and updated_by populated from service layer
	query := `INSERT INTO listing_visits (
		listing_identity_id,
		listing_version,
		user_id,
		owner_user_id,
		scheduled_start,
		scheduled_end,
		duration_minutes,
		status,
		type,
		source,
		realtor_notes,
		owner_notes,
		rejection_reason,
		cancel_reason,
		first_owner_action_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Execute insert via instrumented adapter (metrics + tracing)
	result, err := a.ExecContext(ctx, tx, "insert_visit", query,
		entity.ListingIdentityID,
		entity.ListingVersion,
		entity.RequesterUserID,
		entity.OwnerUserID,
		entity.ScheduledStart,
		entity.ScheduledEnd,
		entity.DurationMinutes,
		entity.Status,
		entity.Type,
		entity.Source,
		entity.RealtorNotes,
		entity.OwnerNotes,
		entity.RejectionReason,
		entity.CancelReason,
		entity.FirstOwnerActionAt,
	)
	if err != nil {
		// Mark span as error for distributed tracing
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.visit.insert.exec_error", "listing_identity_id", entity.ListingIdentityID, "err", err)
		return 0, fmt.Errorf("insert visit: %w", err)
	}

	// Retrieve auto-generated ID from database
	id, err := result.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.visit.insert.last_id_error", "listing_identity_id", entity.ListingIdentityID, "err", err)
		return 0, fmt.Errorf("visit last insert id: %w", err)
	}

	// Update domain object with generated ID for upstream use
	visit.SetID(id)
	return id, nil
}
