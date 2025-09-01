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

func (ua *UserAdapter) UpdateUserByID(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `UPDATE users SET
			full_name = ?, nick_name = ?, national_id = ?, creci_number = ?, creci_state = ?, creci_validity = ?,
			born_at = ?, phone_number = ?, email = ?, zip_code = ?, street = ?, number = ?, complement = ?, neighborhood = ?, 
			city = ?, state = ?, opt_status = ?, deleted = ?, last_signin_attempt = ?
			WHERE id = ?;`

	entity := userconverters.UserDomainToEntity(user)

	_, err = ua.Update(ctx, tx, query,
		entity.FullName,
		entity.NickName,
		entity.NationalID,
		entity.CreciNumber,
		entity.CreciState,
		entity.CreciValidity,
		entity.BornAT,
		entity.PhoneNumber,
		entity.Email,
		entity.ZipCode,
		entity.Street,
		entity.Number,
		entity.Complement,
		entity.Neighborhood,
		entity.City,
		entity.State,
		entity.OptStatus,
		entity.Deleted,
		entity.LastSignInAttempt,
		entity.ID,
	)
	if err != nil {
		slog.Error("mysqluseradapter/UpdateUserByID: error executing Update", "error", err)
		return fmt.Errorf("update user by id: %w", err)
	}

	// Note: User role updates are now handled by permission service

	return
}
