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

func (ua *UserAdapter) CreateAgencyInvite(ctx context.Context, tx *sql.Tx, agency usermodel.UserInterface, phoneNumber string) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	sql := `INSERT INTO agency_invites (agency_id, phone_number) VALUES (?, ?);`

	id, err := ua.Create(ctx, tx, sql, agency.GetID(), phoneNumber)
	if err != nil {
		slog.Error("mysqluseradapter/CreateAgencyInvite: error executing Create", "error", err)
		return status.Error(codes.Internal, "Internal server error")
	}

	agency.SetID(id)

	return
}
