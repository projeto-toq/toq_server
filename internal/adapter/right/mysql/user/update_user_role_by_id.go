package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"

	userconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) UpdateUserRole(ctx context.Context, tx *sql.Tx, role usermodel.UserRoleInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `UPDATE user_roles SET active = ?, status = ?, status_reason = ? WHERE id = ?;`

	entity := userconverters.UserRoleDomainToEntity(role)

	_, err = ua.Update(ctx, tx, query,
		entity.Active,
		entity.Status,
		entity.StatusReason,
		entity.ID,
	)
	if err != nil {
		slog.Error("mysqluseradapter/UpdateUserRole: error executing Update", "error", err)
		err = utils.ErrInternalServer
		return
	}

	return
}
