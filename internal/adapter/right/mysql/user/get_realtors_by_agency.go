package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
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

	entities, err := ua.Read(ctx, tx, query, agencyID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_realtors_by_agency.read_error", "error", err)
		return nil, fmt.Errorf("get realtors by agency read: %w", err)
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
