package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) CreateAgencyInvite(ctx context.Context, tx *sql.Tx, agency usermodel.UserInterface, phoneNumber string) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	sql := `INSERT INTO agency_invites (agency_id, phone_number) VALUES (?, ?);`

	id, err := ua.Create(ctx, tx, sql, agency.GetID(), phoneNumber)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.create_agency_invite.create_error", "error", err)
		return fmt.Errorf("create agency invite: %w", err)
	}

	agency.SetID(id)

	return
}
