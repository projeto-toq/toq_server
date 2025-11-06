package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetInviteByPhoneNumber retrieves an agency invitation by phone number
//
// This function searches the agency_invites table for invitations sent to a specific
// phone number. In the current implementation, if multiple agencies have invited the
// same phone number, only the first invitation is returned.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - phoneNumber: The phone number to search for (E.164 format expected)
//
// Returns:
//   - invite: InviteInterface containing invitation data (id, agency_id, phone_number)
//   - error: sql.ErrNoRows if no invitation found, or database errors
//
// Business Rules:
//   - Phone number is NOT unique in agency_invites (multiple agencies can invite same number)
//   - Returns FIRST matching invitation if multiple exist
//   - Service layer should handle multi-agency invitation logic
//   - Does not validate phone format (service layer responsibility)
//
// Query Logic:
//   - SELECT from agency_invites WHERE phone_number = ?
//   - No JOIN with users table (lightweight query)
//   - Returns raw invitation record
//
// Edge Cases:
//   - No invitation for phone: Returns sql.ErrNoRows
//   - Multiple agencies invited same phone: Returns first result only
//   - Phone number format variations: Exact match required (no normalization)
//   - Invitation already used/expired: Still returned (service validates state)
//
// Performance:
//   - Uses index on agency_invites.phone_number for fast lookup
//   - Lightweight query (no JOINs), suitable for high-frequency checks
//
// Multi-Agency Invitation Handling:
//   - Current behavior: Returns first match
//   - Service layer options:
//   - Accept first invitation and ignore others
//   - Prompt user to choose which agency to join
//   - Implement invitation priority/ranking logic
//
// Use Cases:
//   - Realtor signup flow: Check if phone has pending agency invitation
//   - Invitation validation during onboarding
//   - Agency invite management: Display sent invitations
//
// Example:
//
//	invite, err := adapter.GetInviteByPhoneNumber(ctx, tx, "+5511999999999")
//	if err == sql.ErrNoRows {
//	    // No agency has invited this phone number
//	}
//	// invite contains agency_id and invitation details
func (ua *UserAdapter) GetInviteByPhoneNumber(ctx context.Context, tx *sql.Tx, phoneNumber string) (invite usermodel.InviteInterface, err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query agency invitation by phone number
	// Note: phone_number is indexed but not unique (multiple agencies can invite same number)
	query := `SELECT id, agency_id, phone_number 
	          FROM agency_invites WHERE phone_number = ?;`

	// Execute query using instrumented adapter
	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, phoneNumber)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_invite_by_phone.query_error", "error", queryErr)
		return nil, fmt.Errorf("get invite by phone number query: %w", queryErr)
	}
	defer rows.Close()

	// Scan rows using type-safe function
	entities, err := scanInviteEntities(rows)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_invite_by_phone.scan_error", "error", err)
		return nil, fmt.Errorf("scan invite by phone rows: %w", err)
	}

	// Handle no results: no agency has invited this phone number
	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	// Convert first entity to domain model using type-safe converter
	// Note: If multiple agencies invited same number, returns first match only
	// Service layer should implement logic for multi-agency invitation handling
	invite = userconverters.AgencyInviteEntityToDomainTyped(entities[0])

	return invite, nil
}
