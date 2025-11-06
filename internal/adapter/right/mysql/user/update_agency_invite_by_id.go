package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateAgencyInviteByID updates an existing agency invite record
//
// This function updates the phone number and agency ID for a pending invite.
// Used when an agency needs to modify an invite before it's accepted by a realtor.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (REQUIRED for consistency)
//   - invite: InviteInterface with ID, phone number, and agency ID to update
//
// Returns:
//   - error: sql.ErrNoRows if invite not found, constraint violations, database errors
//
// Business Rules:
//   - ID must be set and > 0 (identifies invite to update)
//   - Invite must exist in agency_invites table
//   - Agency ID must reference valid user with agency role
//   - Phone number updated to new value (may be used to resend invite)
//
// Database Schema:
//   - Table: agency_invites
//   - Primary Key: id
//   - Foreign Key: agency_id -> users.id
//   - Columns: id, agency_id, phone_number
//
// Edge Cases:
//   - Returns sql.ErrNoRows if invite doesn't exist (service maps to 404)
//   - Invalid agency_id triggers foreign key constraint error
//   - Duplicate phone for same agency may cause constraint violation
//
// Performance:
//   - Single-row UPDATE using PRIMARY KEY (very fast)
//   - Foreign key check adds minimal overhead
//
// Important Notes:
//   - Does NOT validate if phone is already used by accepted realtor
//   - Does NOT send new invite SMS (service layer responsibility)
//   - Transaction managed by service layer
//
// Example:
//
//	invite := usermodel.NewInvite()
//	invite.SetID(inviteID)
//	invite.SetAgencyID(agencyID)
//	invite.SetPhoneNumber("+5511999999999")
//
//	err := adapter.UpdateAgencyInviteByID(ctx, tx, invite)
//	if err == sql.ErrNoRows {
//	    // Handle invite not found (404)
//	} else if err != nil {
//	    // Handle infrastructure error
//	}
func (ua *UserAdapter) UpdateAgencyInviteByID(ctx context.Context, tx *sql.Tx, invite usermodel.InviteInterface) (err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Update agency invite by primary key
	// Note: Foreign key constraint ensures agency_id references valid user
	query := `UPDATE agency_invites SET phone_number = ?, agency_id = ? WHERE id = ?`

	// Convert domain model to database entity
	entity := userconverters.AgencyInviteDomainToEntity(invite)

	// Execute update using instrumented adapter (auto-generates metrics + tracing)
	result, execErr := ua.ExecContext(ctx, tx, "update", query,
		entity.PhoneNumber,
		entity.AgencyID,
		entity.ID,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.update_agency_invite.exec_error", "invite_id", entity.ID, "error", execErr)
		return fmt.Errorf("update agency invite: %w", execErr)
	}

	// Check if invite exists and was updated
	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.update_agency_invite.rows_affected_error", "invite_id", entity.ID, "error", rowsErr)
		return fmt.Errorf("agency invite update rows affected: %w", rowsErr)
	}

	// Return sql.ErrNoRows if invite not found (service layer maps to 404)
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
