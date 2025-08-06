package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, status.Error(codes.Internal, "internal server error")
	}

	if len(entities) == 0 {
		return nil, status.Error(codes.NotFound, "Base roles not found")
	}

	for _, entity := range entities {
		id, ok := entity[0].(int64)
		if !ok {
			slog.Error("mysqluseradapter/GetRealtorsByAgency: error converting ID to int64", "value", entity[0])
			return nil, status.Error(codes.Internal, "Internal server error")
		}
		user, err1 := ua.GetUserByID(ctx, tx, id)
		if err1 != nil {
			return nil, err1
		}

		users = append(users, user)
	}

	return

}
