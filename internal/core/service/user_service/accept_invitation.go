package userservices

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) AcceptInvitation(ctx context.Context, userID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		slog.Error("user.accept_invitation.tx_start_error", "err", err)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("user.accept_invitation.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	_, err = us.acceptInvitation(ctx, tx, userID)
	if err != nil {
		return
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		slog.Error("user.accept_invitation.tx_commit_error", "err", err)
		return utils.InternalError("Failed to commit transaction")
	}

	return
}

func (us *userService) acceptInvitation(ctx context.Context, tx *sql.Tx, userID int64) (invitationID int64, err error) {
	//recover the realtor
	realtor, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		return
	}

	//recover the agency invitation
	invite, err := us.repo.GetInviteByPhoneNumber(ctx, tx, realtor.GetPhoneNumber())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, utils.NotFoundError("Invitation")
		}
		return 0, utils.InternalError("Failed to get invitation")
	}

	//recovery the agency inviting the realtor
	agency, err := us.repo.GetUserByID(ctx, tx, invite.GetAgencyID())
	if err != nil {
		return
	}

	//create the agency <-> realtor relationship
	invitationID, err = us.repo.CreateAgencyRelationship(ctx, tx, agency, realtor)
	if err != nil {
		return
	}

	//delete the invitation
	_, err = us.repo.DeleteInviteByID(ctx, tx, invite.GetID())
	if err != nil {
		return
	}

	// Converter RoleInterface para RoleSlug
	activeRole := realtor.GetActiveRole()
	if activeRole == nil || activeRole.GetRole() == nil {
		return 0, utils.InternalError("Active role missing")
	}
	roleSlug := permissionmodel.RoleSlug(activeRole.GetRole().GetSlug())

	status, reason, _, err := us.updateUserStatus(ctx, tx, roleSlug, usermodel.ActionFinishedInviteAccepted)
	if err != nil {
		return
	}

	// TODO: Implementar atualização de status via permission service
	// Por enquanto, apenas registramos o status calculado nos logs
	_ = status
	_ = reason

	// Notificar a imobiliária sobre a aceitação do convite
	notificationService := us.globalService.GetUnifiedNotificationService()
	emailRequest := globalservice.NotificationRequest{
		Type:    globalservice.NotificationTypeEmail,
		To:      agency.GetEmail(),
		Subject: "Convite Aceito - TOQ",
		Body:    fmt.Sprintf("O corretor %s aceitou seu convite para trabalhar com sua imobiliária!", realtor.GetFullName()),
	}

	err = notificationService.SendNotification(ctx, emailRequest)
	if err != nil {
		return
	}

	err = us.repo.UpdateUserByID(ctx, tx, realtor)
	if err != nil {
		return
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableAgencyInvites, "Convite para relacionamento com imobiliária aceito")
	if err != nil {
		return
	}

	return
}
