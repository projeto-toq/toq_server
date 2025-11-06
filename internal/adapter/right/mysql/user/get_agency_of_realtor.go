package mysqluseradapter

import (
	"context"
	"database/sql"
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

	query := `SELECT u.id, u.full_name, u.nick_name, u.national_id, u.creci_number, u.creci_state, u.creci_validity,
	                 u.born_at, u.phone_number, u.email, u.zip_code, u.street, u.number, u.complement,
	                 u.neighborhood, u.city, u.state, u.password, u.opt_status, u.last_activity_at, u.deleted, u.last_signin_attempt
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

	entities, err := scanUserEntities(rows)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_agency_of_realtor.scan_error", "error", err)
		return nil, fmt.Errorf("scan agency rows: %w", err)
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	if len(entities) > 1 {
		errMultiple := fmt.Errorf("multiple agencies found for realtor: %d", realtorID)
		utils.SetSpanError(ctx, errMultiple)
		logger.Error("mysql.user.get_agency_of_realtor.multiple_agencies_error", "realtor_id", realtorID, "error", errMultiple)
		return nil, errMultiple
	}

	agency = userconverters.UserEntityToDomain(entities[0])

	return agency, nil

}
