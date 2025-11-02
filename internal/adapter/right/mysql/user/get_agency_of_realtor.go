package mysqluseradapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) GetAgencyOfRealtor(ctx context.Context, tx *sql.Tx, realtorID int64) (agency usermodel.UserInterface, err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT u.*
				 FROM users u
				 JOIN realtors_agency ra ON u.id = ra.agency_id
				 WHERE ra.realtor_id = ?`

	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, realtorID)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_agency_of_realtor.query_error", "error", queryErr)
		return nil, queryErr
	}
	defer rows.Close()

	entities, err := rowsToEntities(rows)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_agency_of_realtor.rows_to_entities_error", "error", err)
		return nil, fmt.Errorf("scan agency rows: %w", err)
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	if len(entities) > 1 {
		errMultiple := errors.New("multiple agencies found for realtor")
		utils.SetSpanError(ctx, errMultiple)
		logger.Error("mysql.user.get_agency_of_realtor.multiple_agencies_error", "realtor_id", realtorID, "error", errMultiple)
		return nil, errMultiple
	}

	agency, err = userconverters.UserEntityToDomain(entities[0])
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_agency_of_realtor.convert_error", "error", err)
		return nil, fmt.Errorf("convert agency user entity: %w", err)
	}

	return agency, nil

}
