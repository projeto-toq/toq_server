package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetAgencyOfRealtor retrieves the agency associated with a specific realtor
//
// This function performs a JOIN between users and realtors_agency tables to find
// the agency that a realtor belongs to. The relationship is many-to-one: each realtor
// can be associated with at most one agency at a time.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - realtorID: The unique identifier of the realtor
//
// Returns:
//   - agency: UserInterface containing the agency's complete data
//   - error: sql.ErrNoRows if realtor has no agency or database errors
//
// Business Rules:
//   - Uses INNER JOIN on realtors_agency table
//   - Realtor without agency association returns sql.ErrNoRows
//   - Does NOT filter by deleted status (returns even if agency is soft-deleted)
//   - Multiple agencies per realtor is a data integrity violation (logs error)
//
// Query Logic:
//   - JOIN users u with realtors_agency ra on u.id = ra.agency_id
//   - Filter by ra.realtor_id = ?
//   - Returns all user fields for the agency
//
// Edge Cases:
//   - Realtor not in realtors_agency table: Returns sql.ErrNoRows
//   - Multiple agencies for one realtor: Returns error (data integrity issue)
//   - Agency is soft-deleted (deleted=1): Still returned (service decides handling)
//
// Performance:
//   - Uses index on realtors_agency.realtor_id for fast lookup
//   - Single JOIN operation, minimal overhead
//
// Use Cases:
//   - Verifying realtor-agency relationship during authorization
//   - Displaying agency information in realtor profile
//   - Commission calculation workflows
//
// Example:
//
//	agency, err := adapter.GetAgencyOfRealtor(ctx, tx, realtorID)
//	if err == sql.ErrNoRows {
//	    // Realtor is not associated with any agency
//	}
//	// agency contains the associated agency's data
func (ua *UserAdapter) GetAgencyOfRealtor(ctx context.Context, tx *sql.Tx, realtorID int64) (agency usermodel.UserInterface, err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query agency via JOIN with realtors_agency relationship table
	// Note: INNER JOIN excludes realtors without agency association
	query := `SELECT u.id, u.full_name, u.nick_name, u.national_id, u.creci_number, u.creci_state, u.creci_validity,
	                 u.born_at, u.phone_number, u.email, u.zip_code, u.street, u.number, u.complement,
	                 u.neighborhood, u.city, u.state, u.password, u.opt_status, u.last_activity_at, u.deleted, u.last_signin_attempt
				 FROM users u
				 JOIN realtors_agency ra ON u.id = ra.agency_id
				 WHERE ra.realtor_id = ?`

	// Execute query using instrumented adapter
	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, realtorID)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_agency_of_realtor.query_error", "error", queryErr)
		return nil, queryErr
	}
	defer rows.Close()

	// Convert database rows to strongly-typed entities
	entities, err := scanUserEntities(rows)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_agency_of_realtor.scan_error", "error", err)
		return nil, fmt.Errorf("scan agency rows: %w", err)
	}

	// Handle no results: realtor is not associated with any agency
	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	// Safety check: each realtor should have at most one active agency
	// Multiple agencies indicate data integrity violation
	if len(entities) > 1 {
		errMultiple := fmt.Errorf("multiple agencies found for realtor: %d", realtorID)
		utils.SetSpanError(ctx, errMultiple)
		logger.Error("mysql.user.get_agency_of_realtor.multiple_agencies_error", "realtor_id", realtorID, "error", errMultiple)
		return nil, errMultiple
	}

	// Convert database entity to domain model
	agency = userconverters.UserEntityToDomain(entities[0])

	return agency, nil

}
