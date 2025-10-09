package userservices

import (
	"context"
	"database/sql"
	"errors"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (us *userService) GetCrecisToValidateByStatus(ctx context.Context, UserRoleStatus permissionmodel.UserRoleStatus) (realtors []usermodel.UserInterface, err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)

	// Start a database transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.get_crecis_to_validate_by_status.tx_start_error", "error", err)
		return nil, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				utils.LoggerFromContext(ctx).Error("user.get_crecis_to_validate_by_status.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	realtors, err = us.getCrecisToValidateByStatus(ctx, tx, UserRoleStatus)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.get_crecis_to_validate_by_status.tx_commit_error", "error", err)
		return nil, utils.InternalError("Failed to commit transaction")
	}
	return
}

func (us *userService) getCrecisToValidateByStatus(ctx context.Context, tx *sql.Tx, UserRoleStatus permissionmodel.UserRoleStatus) (realtors []usermodel.UserInterface, err error) {
	ctx = utils.ContextWithLogger(ctx)
	// Buscar corretores pelo novo método filtrando por role e status
	realtors, err = us.repo.GetUsersByRoleAndStatus(ctx, tx, permissionmodel.RoleSlugRealtor, UserRoleStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.get_crecis_to_validate_by_status.read_error", "status", UserRoleStatus, "error", err)
		return nil, utils.InternalError("Failed to list realtors by status")
	}
	return realtors, nil
}
