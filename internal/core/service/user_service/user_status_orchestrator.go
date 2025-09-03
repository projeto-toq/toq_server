package userservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	policyport "github.com/giulio-alfieri/toq_server/internal/core/port/policy"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// applyTransition is an internal helper to load context and apply a policy decision.
func (us *userService) applyTransition(ctx context.Context, tx *sql.Tx, action usermodel.ActionFinished) (permissionmodel.UserRoleStatus, globalmodel.NotificationType, error) {
	// Carregar usuário e papel ativo
	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		return 0, 0, utils.ErrInternalServer
	}
	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		return 0, 0, err
	}
	active := user.GetActiveRole()
	if active == nil || active.GetRole() == nil {
		return 0, 0, utils.ErrConflict
	}
	roleSlug := permissionmodel.RoleSlug(active.GetRole().GetSlug())
	from := active.GetStatus()

	// Avaliar política
	to, notif, changed, err := us.statusPolicy.Evaluate(ctx, roleSlug, from, action)
	if err != nil {
		return 0, 0, err
	}
	if !changed {
		return from, notif, nil
	}

	// Persistir alteração
	if err := us.repo.UpdateUserRoleStatus(ctx, tx, userID, roleSlug, to); err != nil {
		return 0, 0, err
	}

	// Auditoria
	if err := us.globalService.CreateAudit(ctx, tx, globalmodel.TableUserRoles, "Atualização de status de role do usuário"); err != nil {
		return 0, 0, err
	}

	return to, notif, nil
}

// ApplyUserStatusTransitionAfterEmailConfirmed decides the next status after email confirmation.
// It checks if phone is still pending and chooses the action accordingly.
func (us *userService) ApplyUserStatusTransitionAfterEmailConfirmed(ctx context.Context) (permissionmodel.UserRoleStatus, globalmodel.NotificationType, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, 0, err
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return 0, 0, err
	}

	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		us.globalService.RollbackTransaction(ctx, tx)
		return 0, 0, utils.ErrInternalServer
	}
	validations, err := us.repo.GetUserValidations(ctx, tx, userID)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return 0, 0, err
	}

	action := usermodel.ActionProfileVerificationCompleted
	if validations.GetPhoneCode() != "" { // telefone ainda pendente
		action = usermodel.ActionProfileEmailVerifiedPhonePending
	}

	to, notif, err := us.applyTransition(ctx, tx, action)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return 0, 0, err
	}

	if err := us.globalService.CommitTransaction(ctx, tx); err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return 0, 0, err
	}

	// Pós-commit: enviar notificação se houver
	if notif != 0 {
		_ = us.globalService.GetUnifiedNotificationService().SendNotification(ctx, globalservice.NotificationRequest{
			Type:    globalservice.NotificationTypeFCM,
			Subject: "Status updated",
			Body:    "Your account status has been updated.",
		})
	}

	return to, notif, nil
}

// ApplyUserStatusTransitionAfterPhoneConfirmed decides the next status after phone confirmation.
// It checks if email is still pending and chooses the action accordingly.
func (us *userService) ApplyUserStatusTransitionAfterPhoneConfirmed(ctx context.Context) (permissionmodel.UserRoleStatus, globalmodel.NotificationType, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, 0, err
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return 0, 0, err
	}

	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		us.globalService.RollbackTransaction(ctx, tx)
		return 0, 0, utils.ErrInternalServer
	}
	validations, err := us.repo.GetUserValidations(ctx, tx, userID)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return 0, 0, err
	}

	action := usermodel.ActionProfileVerificationCompleted
	if validations.GetEmailCode() != "" { // email ainda pendente
		action = usermodel.ActionProfilePhoneVerifiedEmailPending
	}

	to, notif, err := us.applyTransition(ctx, tx, action)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return 0, 0, err
	}

	if err := us.globalService.CommitTransaction(ctx, tx); err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return 0, 0, err
	}

	if notif != 0 {
		_ = us.globalService.GetUnifiedNotificationService().SendNotification(ctx, globalservice.NotificationRequest{
			Type:    globalservice.NotificationTypeFCM,
			Subject: "Status updated",
			Body:    "Your account status has been updated.",
		})
	}

	return to, notif, nil
}

// statusPolicy is injected via factory/bootstrap; define a small interface contract usage here.
// We use the full port interface to preserve compile-time checks.
var _ policyport.UserStatusPolicy
