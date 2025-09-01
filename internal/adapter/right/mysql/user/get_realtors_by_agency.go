package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) GetRealtorsByAgency(ctx context.Context, tx *sql.Tx, agencyID int64) (users []usermodel.UserInterface, err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT realtor_id from realtors_agency WHERE agency_id = ?;`

	entities, err := ua.Read(ctx, tx, query, agencyID)
	if err != nil {
		slog.Error("mysqluseradapter/GetRealtorsByAgency: error executing Read", "error", err)
		return nil, fmt.Errorf("get realtors by agency read: %w", err)
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	for _, entity := range entities {
		id, ok := entity[0].(int64)
		if !ok {
			slog.Error("mysqluseradapter/GetRealtorsByAgency: error converting ID to int64", "value", entity[0])
			return nil, fmt.Errorf("get realtors by agency: invalid id type %T", entity[0])
		}
		user, err1 := ua.GetUserByID(ctx, tx, id)
		if err1 != nil {
			return nil, err1
		}

		users = append(users, user)
	}

	return

}
