package userservices

import (
	"context"
	"database/sql"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) GetOnboardingStatus(ctx context.Context, userID int64) (UserRoleStatus string, reason string, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Start transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	UserRoleStatus, reason, err = us.getOnboardingStatus(ctx, tx, userID)
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

func (us *userService) getOnboardingStatus(ctx context.Context, tx *sql.Tx, userID int64) (UserRoleStatus string, reason string, err error) {

	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		return
	}

	switch user.GetActiveRole().GetStatus() {
	case usermodel.StatusActive:
		UserRoleStatus = "Onboarding finished. User is active."
		reason = user.GetActiveRole().GetStatusReason()
	case usermodel.StatusBlocked:
		UserRoleStatus = "User is blocked."
		reason = user.GetActiveRole().GetStatusReason()
	case usermodel.StatusPendingProfile:
		UserRoleStatus = "System is waiting for user to confirm phone and/or email."
		reason = user.GetActiveRole().GetStatusReason()
	case usermodel.StatusPendingImages:
		UserRoleStatus = "System is waiting for user to upload creci images and selfie."
		reason = user.GetActiveRole().GetStatusReason()
	case usermodel.StatusPendingOCR:
		UserRoleStatus = "System is waiting for AI verification of creci images."
		reason = user.GetActiveRole().GetStatusReason()
	case usermodel.StatusRejectByOCR:
		UserRoleStatus = "AI rejected creci images."
		reason = user.GetActiveRole().GetStatusReason()
	case usermodel.StatusPendingFace:
		UserRoleStatus = "System is waiting for AI verification of user's selfie."
		reason = user.GetActiveRole().GetStatusReason()
	case usermodel.StatusRejectByFace:
		UserRoleStatus = "AI rejected user selfie."
		reason = user.GetActiveRole().GetStatusReason()
	case usermodel.StatusPendingManual:
		UserRoleStatus = "System is waiting for manual verification."
		reason = user.GetActiveRole().GetStatusReason()
	case usermodel.StatusInvitePending:
		UserRoleStatus = "User is invited to join an agency team."
		reason = user.GetActiveRole().GetStatusReason()
	}
	return
}
