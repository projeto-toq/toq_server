package userservices

import (
	"context"
	"database/sql"
	"log/slog"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) GetCrecisToValidateByStatus(ctx context.Context, UserRoleStatus permissionmodel.UserRoleStatus) (realtors []usermodel.UserInterface, err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Start a database transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	realtors, err = us.getCrecisToValidateByStatus(ctx, tx, UserRoleStatus)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	// Commit the transaction
	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}
	return
}

func (us *userService) getCrecisToValidateByStatus(_ context.Context, _ *sql.Tx, UserRoleStatus permissionmodel.UserRoleStatus) (realtors []usermodel.UserInterface, err error) {

	// TODO: Reimplementar busca por status após migração do sistema de status
	// O método GetUsersByStatus precisa ser atualizado para o novo sistema de permissões
	// Por enquanto, retornar lista vazia
	slog.Warn("getCrecisToValidateByStatus temporarily disabled during migration", "status", UserRoleStatus)
	return []usermodel.UserInterface{}, nil

	/*
		// Código original comentado durante migração:
		// Read the realtors user with given status from the database
		realtors, err = us.repo.GetUsersByStatus(ctx, tx, UserRoleStatus, permissionmodel.RoleSlugRealtor)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				slog.Error("Failed to read realtor users with status in GetCrecisToValidateByStatus", "error", err)
				return
			}
			return nil, nil
		}

		return
	*/
}
