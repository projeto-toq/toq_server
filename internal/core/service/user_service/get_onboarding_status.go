package userservices

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) GetOnboardingStatus(ctx context.Context, userID int64) (UserRoleStatus string, reason string, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return "", "", utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	// Start transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		slog.Error("user.get_onboarding_status.tx_start_error", "err", err)
		return "", "", utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("user.get_onboarding_status.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	UserRoleStatus, reason, err = us.getOnboardingStatus(ctx, tx, userID)
	if err != nil {
		return
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		slog.Error("user.get_onboarding_status.tx_commit_error", "err", err)
		return "", "", utils.InternalError("Failed to commit transaction")
	}

	return
}

func (us *userService) getOnboardingStatus(ctx context.Context, tx *sql.Tx, userID int64) (UserRoleStatus string, reason string, err error) {

	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		return
	}

	// TODO: Reimplementar verificação de status após migração completa do sistema de status
	// Por enquanto, retornar status genérico baseado na existência de role ativo
	if user.GetActiveRole() != nil {
		UserRoleStatus = "User has active role - status system under migration"
		reason = "Status verification temporarily disabled during migration"
	} else {
		UserRoleStatus = "No active role found"
		reason = "User has no active role assigned"
	}
	return

	/* Código original comentado durante migração:
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
	*/
}
