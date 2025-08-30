package userservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) AddAlternativeRole(ctx context.Context, userID int64, roleSlug permissionmodel.RoleSlug, creciInfo ...string) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	err = us.addAlternativeRole(ctx, tx, userID, roleSlug, creciInfo...)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	return
}

func (us *userService) addAlternativeRole(ctx context.Context, tx *sql.Tx, userID int64, roleSlug permissionmodel.RoleSlug, creciInfo ...string) (err error) {

	//verify if the user is on active status
	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		return
	}

	// Check if user has active role
	activeRole := user.GetActiveRole()
	if activeRole == nil {
		err = utils.ErrInternalServer
		return
	}

	// Validate creci info for realtor role
	if roleSlug == permissionmodel.RoleSlugRealtor && len(creciInfo) != 3 {
		err = utils.ErrInternalServer
		return
	}

	// Get role from permission service
	role, err := us.permissionService.GetRoleBySlugWithTx(ctx, tx, roleSlug)
	if err != nil {
		return
	}

	// Create user role using permission service (not active by default)
	err = us.permissionService.AssignRoleToUserWithTx(ctx, tx, userID, role.GetID(), nil)
	if err != nil {
		return
	}

	// Handle realtor-specific setup
	if roleSlug == permissionmodel.RoleSlugRealtor {
		err = us.CreateUserFolder(ctx, user.GetID())
		if err != nil {
			return
		}
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUserRoles, "Criado papel alternativo")
	if err != nil {
		return
	}

	return
}
