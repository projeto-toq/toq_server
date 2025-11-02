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

	query := `SELECT realtor_id from realtors_agency WHERE agency_id = ?;`

	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, agencyID)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_realtors_by_agency.query_error", "error", queryErr)
		return nil, fmt.Errorf("get realtors by agency query: %w", queryErr)
	}
	defer rows.Close()

	entities, err := rowsToEntities(rows)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_realtors_by_agency.rows_to_entities_error", "error", err)
		return nil, fmt.Errorf("scan realtors by agency rows: %w", err)
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	for _, entity := range entities {
		id, ok := entity[0].(int64)
		if !ok {
			errInvalid := fmt.Errorf("get realtors by agency: invalid id type %T", entity[0])
			utils.SetSpanError(ctx, errInvalid)
			logger.Error("mysql.user.get_realtors_by_agency.invalid_id_type", "value", entity[0], "error", errInvalid)
			return nil, errInvalid
		}
		user, err1 := ua.GetUserByID(ctx, tx, id)
		if err1 != nil {
			utils.SetSpanError(ctx, err1)
			logger.Error("mysql.user.get_realtors_by_agency.get_user_error", "user_id", id, "error", err1)
			return nil, fmt.Errorf("get realtor by id: %w", err1)
		}

		users = append(users, user)
	}

	return

}
