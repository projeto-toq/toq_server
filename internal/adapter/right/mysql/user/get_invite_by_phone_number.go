package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) GetInviteByPhoneNumber(ctx context.Context, tx *sql.Tx, phoneNumber string) (invite usermodel.InviteInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT * FROM agency_invites WHERE phone_number = ?;`

	entities, err := ua.Read(ctx, tx, query, phoneNumber)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_invite_by_phone.read_error", "error", err)
		return nil, fmt.Errorf("get invite by phone number read: %w", err)
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	invite, err = userconverters.AgencyInviteEntityToDomain(entities[0])
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_invite_by_phone.convert_error", "error", err)
		return nil, fmt.Errorf("convert agency invite entity: %w", err)
	}

	return invite, nil
}
