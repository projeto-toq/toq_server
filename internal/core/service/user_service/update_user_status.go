package userservices

import (
	"context"
	"database/sql"
	"errors"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

func (us *userService) updateUserStatus(
	ctx context.Context,
	tx *sql.Tx,
	role permissionmodel.RoleSlug,
	actonFinished usermodel.ActionFinished,
	user ...usermodel.UserInterface) (
	nextStatus permissionmodel.UserRoleStatus,
	nextStatusReason string,
	notification globalmodel.NotificationType,
	err error) {

	// No tracer here: private helper should not create spans; caller's span will capture infra errors.

	switch actonFinished {
	case usermodel.ActionFinishedCreated: //after user is created
		switch role {
		case permissionmodel.RoleSlugRoot:
			nextStatus = permissionmodel.StatusActive
			nextStatusReason = "Admin user"
		default:
			nextStatus = permissionmodel.StatusPendingEmail
			nextStatusReason = "Phone and email validation pending"
		}
	case usermodel.ActionFinishedPhoneVerified: //after phone is verified, what is always first
		nextStatus = permissionmodel.StatusPendingEmail
		nextStatusReason = "Email validation pending"
	case usermodel.ActionFinishedEmailVerified: //after email is verified. the profile is OK
		switch role {
		case permissionmodel.RoleSlugOwner:
			nextStatus = permissionmodel.StatusActive
			nextStatusReason = "User active"
		// case permissionmodel.RoleSlugRealtor:
		// 	nextStatus = permissionmodel.StatusPendingImages
		// 	nextStatusReason = "Awaiting creci images to verify"
		case permissionmodel.RoleSlugAgency:
			nextStatus = permissionmodel.StatusPendingManual
			nextStatusReason = "Awaiting administrator approval"
		}
	// case usermodel.ActionFinishedCreciImagesUploaded:
	// 	nextStatus = permissionmodel.StatusPendingOCR
	// 	nextStatusReason = "Awaiting OCR verification"
	// case usermodel.ActionFinishedCreciNumberDoesntMatch:
	// 	nextStatus = permissionmodel.StatusRejectByOCR
	// 	nextStatusReason = "Creci number doesn't match"
	// 	notification = globalmodel.NotificationInvalidCreciNumber
	// case usermodel.ActionFinishedCreciStateDoesntMatch:
	// 	nextStatus = permissionmodel.StatusRejectByOCR
	// 	nextStatusReason = "Creci state doesn't match"
	// 	notification = globalmodel.NotificationInvalidCreciState
	// case usermodel.ActionFinishedCreciStateUnsupported:
	// 	nextStatus = permissionmodel.StatusRejectByOCR
	// 	nextStatusReason = "Creci state unsupported"
	// 	notification = globalmodel.NotificationCreciStateUnsupported
	// case usermodel.ActionFinishedBadCreciImages:
	// 	nextStatus = permissionmodel.StatusRejectByOCR
	// 	nextStatusReason = "Bad creci images"
	// 	notification = globalmodel.NotificationBadCreciImages
	// case usermodel.ActionFinishedBadSelfieImage:
	// 	nextStatus = permissionmodel.StatusRejectByFace
	// 	nextStatusReason = "Bad Selfie image"
	// 	notification = globalmodel.NotificationBadSelfieImage
	// case usermodel.ActionFinishedSelfieDoesntMatch:
	// 	nextStatus = permissionmodel.StatusRejectByFace
	// 	nextStatusReason = "Selfie doesn't match"
	// 	notification = globalmodel.NotificationBadSelfieImage
	// case usermodel.ActionFinishedCreciVerified:
	// 	nextStatus = permissionmodel.StatusPendingFace
	// 	nextStatusReason = "Awaiting face verification"
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
			nextStatus = permissionmodel.StatusInvitePending
			nextStatusReason = "Awaiting invite verification"
			notification = globalmodel.NotificationRealtorInvitePush
		} else {
			nextStatus = permissionmodel.StatusActive
			nextStatusReason = "User active"
			notification = globalmodel.NotificationCreciValidated
		}

	case usermodel.ActionFinishedCreciManualVerified:
		nextStatus = permissionmodel.StatusActive
		nextStatusReason = "User active"
		notification = globalmodel.NotificationCreciValidated
	case usermodel.ActionFinishedCreciImagesUploadedForManualReview:
		nextStatus = permissionmodel.StatusPendingManual
		nextStatusReason = "Awaiting manual verification by administrator"
	case usermodel.ActionFinishedInviteCreated:
		nextStatus = permissionmodel.StatusInvitePending
		nextStatusReason = "Awaiting invite verification"
		notification = globalmodel.NotificationRealtorInvitePush
	case usermodel.ActionFinishedInviteAccepted:
		nextStatus = permissionmodel.StatusActive
		nextStatusReason = "User active"
		notification = globalmodel.NotificationInviteAccepted
	case usermodel.ActionFinishedInviteRejected:
		nextStatus = permissionmodel.StatusActive
		nextStatusReason = "User active"
		notification = globalmodel.NotificationInviteRejected
	}

	return
}
