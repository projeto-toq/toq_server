package userservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (us *userService) AcceptInvitation(ctx context.Context, userID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	_, err = us.acceptInvitation(ctx, tx, userID)
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

func (us *userService) acceptInvitation(ctx context.Context, tx *sql.Tx, userID int64) (invitationID int64, err error) {
	//recover the realtor
	realtor, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		return
	}

	//recover the agency invitation
	invite, err := us.repo.GetInviteByPhoneNumber(ctx, tx, realtor.GetPhoneNumber())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return 0, status.Error(codes.FailedPrecondition, "Invite not found")
		}
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

	status, reason, notification, err := us.updateUserStatus(ctx, tx, realtor.GetActiveRole().GetRole(), usermodel.ActionFinishedInviteAccepted)
	if err != nil {
		return
	}
	realtor.GetActiveRole().SetStatus(status)
	realtor.GetActiveRole().SetStatusReason(reason)
	err = us.globalService.SendNotification(ctx, agency, notification, realtor.GetFullName())
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
