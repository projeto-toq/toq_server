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

func (ua *UserAdapter) GetInviteByPhoneNumber(ctx context.Context, tx *sql.Tx, phoneNumber string) (invite usermodel.InviteInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT * FROM agency_invites WHERE phone_number = ?;`

	entities, err := ua.Read(ctx, tx, query, phoneNumber)
	if err != nil {
		slog.Error("mysqluseradapter/GetInviteByPhoneNumber: error executing Read", "error", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	if len(entities) == 0 {
		return nil, status.Error(codes.NotFound, "Invite not found")
	}

	invite, err = userconverters.AgencyInviteEntityToDomain(entities[0])
	if err != nil {
		return
	}

	return invite, nil
}
