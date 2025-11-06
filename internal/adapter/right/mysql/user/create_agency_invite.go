package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateAgencyInvite creates an invitation for a realtor to join an agency
//
// This function inserts a new record in the agency_invites table, allowing
// an agency to send invitations to realtors via phone number. The realtor
// can later accept the invite during onboarding.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (REQUIRED - invite creation must be transactional)
//   - agency: UserInterface representing the agency sending the invitation
//   - phoneNumber: Target realtor's phone number in E.164 format
//
// Returns:
//   - error: Database errors (duplicate invite, foreign key violations)
//
// Side Effects:
//   - ⚠️ MODIFIES agency.ID with the invite ID (potential bug - should not modify input)
//
// Business Rules:
//   - Agency must exist (foreign key constraint on agency_id)
//   - Phone number must be valid E.164 format (enforced by service layer)
//   - Duplicate invites for same agency+phone are allowed (no unique constraint)
//
// Database Constraints:
//   - FK: agency_id REFERENCES users(id)
//
// Note: This function has a side effect of modifying agency.ID with the invite ID.
// Consider refactoring to return inviteID instead of modifying input parameter.
func (ua *UserAdapter) CreateAgencyInvite(ctx context.Context, tx *sql.Tx, agency usermodel.UserInterface, phoneNumber string) (err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Insert new agency invite record
	sql := `INSERT INTO agency_invites (agency_id, phone_number) VALUES (?, ?);`

	// Execute insert using instrumented adapter (auto-generates metrics + tracing)
	result, execErr := ua.ExecContext(ctx, tx, "insert", sql, agency.GetID(), phoneNumber)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.create_agency_invite.exec_error", "error", execErr)
		return fmt.Errorf("create agency invite: %w", execErr)
	}

	// Retrieve auto-generated primary key
	id, lastErr := result.LastInsertId()
	if lastErr != nil {
		utils.SetSpanError(ctx, lastErr)
		logger.Error("mysql.user.create_agency_invite.last_insert_id_error", "error", lastErr)
		return fmt.Errorf("agency invite last insert id: %w", lastErr)
	}

	// ⚠️ Side effect: modifies input parameter (consider refactoring)
	agency.SetID(id)

	return
}
