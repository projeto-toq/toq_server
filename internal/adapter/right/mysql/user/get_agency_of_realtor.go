package mysqluseradapter

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	userconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) GetAgencyOfRealtor(ctx context.Context, tx *sql.Tx, realtorID int64) (agency usermodel.UserInterface, err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT u.*
				 FROM users u
				 JOIN realtors_agency ra ON u.id = ra.agency_id
				 WHERE ra.realtor_id = ?`

	entities, err := ua.Read(ctx, tx, query, realtorID)
	if err != nil {
		slog.Error("mysqluseradapter/GetAgencyOfRealtor: error executing Read", "error", err)
		return nil, err
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	if len(entities) > 1 {
		slog.Error("mysqluseradapter.GetAgencyOfRealtor: Multiple agencies found for the same realtor ID", "realtorID", realtorID)
		return nil, errors.New("multiple agencies found for realtor")
	}

	return userconverters.UserEntityToDomain(entities[0])

}
