package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) UpdateUserByID(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE users SET
			full_name = ?, nick_name = ?, national_id = ?, creci_number = ?, creci_state = ?, creci_validity = ?,
			born_at = ?, phone_number = ?, email = ?, zip_code = ?, street = ?, number = ?, complement = ?, neighborhood = ?, 
			city = ?, state = ?, opt_status = ?, deleted = ?, last_signin_attempt = ?
			WHERE id = ?;`

	entity := userconverters.UserDomainToEntity(user)

	result, execErr := ua.ExecContext(ctx, tx, "update", query,
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
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.update_user_by_id.exec_error", "error", execErr)
		return fmt.Errorf("update user by id: %w", execErr)
	}

	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.update_user_by_id.rows_affected_error", "error", rowsErr)
		return fmt.Errorf("user update rows affected: %w", rowsErr)
	}

	// Note: User role updates are now handled by permission service

	return
}
