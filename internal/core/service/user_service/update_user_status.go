package userservices

import (
	"context"
	"database/sql"
	"errors"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

func (us *userService) updateUserStatus(
	ctx context.Context,
	tx *sql.Tx,
	role permissionmodel.RoleSlug,
	actonFinished usermodel.ActionFinished,
	user ...usermodel.UserInterface) (
	nextStatus globalmodel.UserRoleStatus,
	nextStatusReason string,
	notification globalmodel.NotificationType,
	err error) {

	// No tracer here: private helper should not create spans; caller's span will capture infra errors.

	switch actonFinished {
	case usermodel.ActionFinishedCreated: //after user is created
		switch role {
		case permissionmodel.RoleSlugRoot:
			nextStatus = globalmodel.StatusActive
			nextStatusReason = "Admin user"
		default:
			nextStatus = globalmodel.StatusPendingEmail
			nextStatusReason = "Phone and email validation pending"
		}
	case usermodel.ActionFinishedPhoneVerified: //after phone is verified, what is always first
		nextStatus = globalmodel.StatusPendingEmail
		nextStatusReason = "Email validation pending"
	case usermodel.ActionFinishedEmailVerified: //after email is verified. the profile is OK
		switch role {
		case permissionmodel.RoleSlugOwner:
			nextStatus = globalmodel.StatusActive
			nextStatusReason = "User active"
		// case permissionmodel.RoleSlugRealtor:
		// 	nextStatus = permissionmodel.StatusPendingImages
		// 	nextStatusReason = "Awaiting creci images to verify"
		case permissionmodel.RoleSlugAgency:
			nextStatus = globalmodel.StatusPendingManual
			nextStatusReason = "Awaiting administrator approval"
		}
	case usermodel.ActionFinishedCreciFaceVerified:
		iUser := usermodel.NewUser()
		if len(user) > 0 {
			iUser = user[0]
		}
		exist := false
		_, err := us.repo.GetInviteByPhoneNumber(ctx, tx, iUser.GetPhoneNumber())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return 0, "", 0, err
			}
		} else {
			exist = true
		}

		if exist {
			nextStatus = globalmodel.StatusPendingManual // StatusInvitePending removido: reutiliza pending_manual para aguardar aceite do convite
			nextStatusReason = "Awaiting invite verification"
			notification = globalmodel.NotificationRealtorInvitePush
		} else {
			nextStatus = globalmodel.StatusActive
			nextStatusReason = "User active"
			notification = globalmodel.NotificationCreciValidated
		}

	case usermodel.ActionFinishedCreciManualVerified:
		nextStatus = globalmodel.StatusActive
		nextStatusReason = "User active"
		notification = globalmodel.NotificationCreciValidated
	case usermodel.ActionFinishedCreciImagesUploadedForManualReview:
		nextStatus = globalmodel.StatusPendingManual
		nextStatusReason = "Awaiting manual verification by administrator"
	case usermodel.ActionFinishedInviteCreated:
		nextStatus = globalmodel.StatusPendingManual // StatusInvitePending removido: mantém usuário aguardando aceite manual
		nextStatusReason = "Awaiting invite verification"
		notification = globalmodel.NotificationRealtorInvitePush
	case usermodel.ActionFinishedInviteAccepted:
		nextStatus = globalmodel.StatusActive
		nextStatusReason = "User active"
		notification = globalmodel.NotificationInviteAccepted
	case usermodel.ActionFinishedInviteRejected:
		nextStatus = globalmodel.StatusActive
		nextStatusReason = "User active"
		notification = globalmodel.NotificationInviteRejected
	}

	return
}
