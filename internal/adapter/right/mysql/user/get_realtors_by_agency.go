package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetRealtorsByAgency retrieves all realtors associated with a specific agency
//
// This function performs a two-step query to fetch realtor data:
//  1. Query realtors_agency table for realtor IDs associated with the agency
//  2. Fetch full user data for each realtor ID via GetUserByID
//
// This N+1 query pattern is intentional to reuse existing GetUserByID logic
// and maintain consistency in user data retrieval across the application.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - agencyID: The unique identifier of the agency
//
// Returns:
//   - users: Slice of UserInterface containing all realtors for the agency
//   - error: sql.ErrNoRows if agency has no realtors, or database errors
//
// Business Rules:
//   - Uses realtors_agency many-to-many relationship table
//   - Returns sql.ErrNoRows if agency has no associated realtors
//   - Includes soft-deleted realtors (GetUserByID filters deleted=0)
//   - Order is not guaranteed (database insertion order)
//
// Query Logic:
//   - Step 1: SELECT realtor_id FROM realtors_agency WHERE agency_id = ?
//   - Step 2: For each ID, call GetUserByID(ctx, tx, realtorID)
//
// Performance Considerations:
//   - N+1 query pattern: One query for IDs, then N queries for user data
//   - Trade-off: Consistency vs performance
//   - For large agencies (>50 realtors), consider batch query optimization
//   - Typical use case: Small to medium agencies (5-20 realtors)
//
// Edge Cases:
//   - Agency with no realtors: Returns sql.ErrNoRows
//   - Realtor ID in realtors_agency but user deleted: GetUserByID returns sql.ErrNoRows, propagated as error
//   - Orphaned realtor_id (user doesn't exist): GetUserByID returns sql.ErrNoRows, propagated as error
//   - Agency is soft-deleted: Still executes query (no filter on agency deleted status)
//
// Alternative Approaches Considered:
//   - Single JOIN query: Would duplicate GetUserByID logic and entity scanning
//   - Batch query with IN clause: Would require separate entity scanning logic
//   - Current approach prioritizes code reuse and maintainability
//
// Use Cases:
//   - Listing all realtors in agency dashboard
//   - Calculating agency-wide statistics (total realtors, active listings, etc.)
//   - Bulk notifications to all realtors in an agency
//   - Agency performance reports
//
// Example:
//
//	realtors, err := adapter.GetRealtorsByAgency(ctx, tx, agencyID)
//	if err == sql.ErrNoRows {
//	    // Agency has no associated realtors
//	}
//	// realtors contains full user data for each realtor
func (ua *UserAdapter) GetRealtorsByAgency(ctx context.Context, tx *sql.Tx, agencyID int64) (users []usermodel.UserInterface, err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Step 1: Query realtor IDs associated with agency
	// Note: realtors_agency is a many-to-many relationship table (agency_id + realtor_id)
	// Returns only realtor_id column, full user data fetched via GetUserByID in step 2
	query := `SELECT realtor_id FROM realtors_agency WHERE agency_id = ?;`

	// Execute query using instrumented adapter
	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, agencyID)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_realtors_by_agency.query_error", "error", queryErr)
		return nil, fmt.Errorf("get realtors by agency query: %w", queryErr)
	}
	defer rows.Close()

	// Step 2: Scan realtor IDs into slice (no entity needed for single column)
	var realtorIDs []int64
	for rows.Next() {
		var realtorID int64
		if err := rows.Scan(&realtorID); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.user.get_realtors_by_agency.scan_error", "error", err)
			return nil, fmt.Errorf("scan realtor_id: %w", err)
		}
		realtorIDs = append(realtorIDs, realtorID)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_realtors_by_agency.rows_error", "error", err)
		return nil, fmt.Errorf("iterate realtor rows: %w", err)
	}

	// Handle no results: agency has no associated realtors
	if len(realtorIDs) == 0 {
		return nil, sql.ErrNoRows
	}

	// Step 3: Fetch full user data for each realtor ID
	// Note: N+1 query pattern for code reuse and consistency
	for _, realtorID := range realtorIDs {
		user, err1 := ua.GetUserByID(ctx, tx, realtorID)
		if err1 != nil {
			utils.SetSpanError(ctx, err1)
			logger.Error("mysql.user.get_realtors_by_agency.get_user_error", "user_id", realtorID, "error", err1)
			return nil, fmt.Errorf("get realtor by id: %w", err1)
		}

		users = append(users, user)
	}

	return

}
