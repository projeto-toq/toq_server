package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) UpdateAgencyInviteByID(ctx context.Context, tx *sql.Tx, invite usermodel.InviteInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE agency_invites SET phone_number = ?, agency_id = ? WHERE id = ?;`

	entity := userconverters.AgencyInviteDomainToEntity(invite)

	result, execErr := ua.ExecContext(ctx, tx, "update", query,
		entity.PhoneNumber,
		entity.AgencyID,
		entity.ID,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.update_agency_invite.exec_error", "error", execErr)
		return fmt.Errorf("update agency invite: %w", execErr)
	}

	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.update_agency_invite.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("agency invite update rows affected: %w", rowsErr)
	}

	return
}
