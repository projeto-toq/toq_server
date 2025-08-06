package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"

	userconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, status.Error(codes.Internal, "internal server error")
	}

	if len(entities) == 0 {
		return nil, status.Error(codes.NotFound, "Agency not found")
	}

	if len(entities) > 1 {
		slog.Error("mysqluseradapter.GetAgencyOfRealtor: Multiple agencies found for the same realtor ID", "realtorID", realtorID)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return userconverters.UserEntityToDomain(entities[0])

}
