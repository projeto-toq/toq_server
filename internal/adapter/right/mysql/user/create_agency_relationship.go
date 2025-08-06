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

func (ua *UserAdapter) CreateAgencyRelationship(ctx context.Context, tx *sql.Tx, agency usermodel.UserInterface, realtor usermodel.UserInterface) (id int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	sql := `INSERT INTO realtors_agency (agency_id, realtor_id) VALUES (?, ?);`

	id, err = ua.Create(ctx, tx, sql, agency.GetID(), realtor.GetID())
	if err != nil {
		slog.Error("mysqluseradapter/CreateAgencyRelationship: error executing Create", "error", err)
		return 0, status.Error(codes.Internal, "Internal server error")
	}

	return id, nil

}
