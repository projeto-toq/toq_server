package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) GetRealtorsByAgency(ctx context.Context, tx *sql.Tx, agencyID int64) (users []usermodel.UserInterface, err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query realtor IDs associated with agency
	// Note: realtors_agency is a many-to-many relationship table (agency_id + realtor_id)
	// Returns only realtor_id column, then fetches full user data via GetUserByID
	query := `SELECT realtor_id FROM realtors_agency WHERE agency_id = ?;`

	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, agencyID)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_realtors_by_agency.query_error", "error", queryErr)
		return nil, fmt.Errorf("get realtors by agency query: %w", queryErr)
	}
	defer rows.Close()

	// Scan realtor IDs directly into int64 slice (no entity needed for single column)
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

	// Handle no results
	if len(realtorIDs) == 0 {
		return nil, sql.ErrNoRows
	}

	// Fetch full user data for each realtor ID
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
