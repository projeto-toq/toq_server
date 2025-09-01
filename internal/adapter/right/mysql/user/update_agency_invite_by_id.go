package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	userconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) UpdateAgencyInviteByID(ctx context.Context, tx *sql.Tx, invite usermodel.InviteInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `UPDATE agency_invites SET phone_number = ?, agency_id = ? WHERE id = ?;`

	entity := userconverters.AgencyInviteDomainToEntity(invite)

	_, err = ua.Update(ctx, tx, query,
		entity.PhoneNumber,
		entity.AgencyID,
		entity.ID,
	)
	if err != nil {
		slog.Error("mysqluseradapter/UpdateAgencyInviteByID: error executing Update", "error", err)
		return fmt.Errorf("update agency invite: %w", err)
	}

	return
}
