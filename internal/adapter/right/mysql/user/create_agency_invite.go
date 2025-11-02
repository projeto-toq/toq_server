package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
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

	result, execErr := ua.ExecContext(ctx, tx, "insert", sql, agency.GetID(), phoneNumber)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.create_agency_invite.exec_error", "error", execErr)
		return fmt.Errorf("create agency invite: %w", execErr)
	}

	id, lastErr := result.LastInsertId()
	if lastErr != nil {
		utils.SetSpanError(ctx, lastErr)
		logger.Error("mysql.user.create_agency_invite.last_insert_id_error", "error", lastErr)
		return fmt.Errorf("agency invite last insert id: %w", lastErr)
	}

	agency.SetID(id)

	return
}
