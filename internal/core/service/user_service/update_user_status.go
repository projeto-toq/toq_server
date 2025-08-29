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

func (us *userService) updateUserStatus(
	ctx context.Context,
	tx *sql.Tx,
	role usermodel.UserRole,
	actonFinished usermodel.ActionFinished,
	user ...usermodel.UserInterface) (
	nextStatus usermodel.UserRoleStatus,
	nextStatusReason string,
	notification globalmodel.NotificationType,
	err error) {

	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	switch actonFinished {
	case usermodel.ActionFinishedCreated: //after user is created
		switch role {
		case usermodel.RoleRoot:
			nextStatus = usermodel.StatusActive
			nextStatusReason = "Admin user"
		default:
			nextStatus = usermodel.StatusPendingProfile
			nextStatusReason = "Phone and email validation pending"
		}
	case usermodel.ActionFinishedPhoneVerified: //after phone is verified, what is always first
		nextStatus = usermodel.StatusPendingProfile
		nextStatusReason = "Email validation pending"
	case usermodel.ActionFinishedEmailVerified: //after email is verified. the profile is OK
		switch role {
		case usermodel.RoleOwner:
			nextStatus = usermodel.StatusActive
			nextStatusReason = "User active"
		case usermodel.RoleRealtor:
			nextStatus = usermodel.StatusPendingImages
			nextStatusReason = "Awaiting creci images to verify"
		case usermodel.RoleAgency:
			nextStatus = usermodel.StatusPendingManual
			nextStatusReason = "Awaiting administrator approval"
		}
	case usermodel.ActionFinishedCreciImagesUploaded:
		nextStatus = usermodel.StatusPendingOCR
		nextStatusReason = "Awaiting OCR verification"
	case usermodel.ActionFinishedCreciNumberDoesntMatch:
		nextStatus = usermodel.StatusRejectByOCR
		nextStatusReason = "Creci number doesn't match"
		notification = globalmodel.NotificationInvalidCreciNumber
	case usermodel.ActionFinishedCreciStateDoesntMatch:
		nextStatus = usermodel.StatusRejectByOCR
		nextStatusReason = "Creci state doesn't match"
		notification = globalmodel.NotificationInvalidCreciState
	case usermodel.ActionFinishedCreciStateUnsupported:
		nextStatus = usermodel.StatusRejectByOCR
		nextStatusReason = "Creci state unsupported"
		notification = globalmodel.NotificationCreciStateUnsupported
	case usermodel.ActionFinishedBadCreciImages:
		nextStatus = usermodel.StatusRejectByOCR
		nextStatusReason = "Bad creci images"
		notification = globalmodel.NotificationBadCreciImages
	case usermodel.ActionFinishedBadSelfieImage:
		nextStatus = usermodel.StatusRejectByFace
		nextStatusReason = "Bad Selfie image"
		notification = globalmodel.NotificationBadSelfieImage
	case usermodel.ActionFinishedSelfieDoesntMatch:
		nextStatus = usermodel.StatusRejectByFace
		nextStatusReason = "Selfie doesn't match"
		notification = globalmodel.NotificationBadSelfieImage
	case usermodel.ActionFinishedCreciVerified:
		nextStatus = usermodel.StatusPendingFace
		nextStatusReason = "Awaiting face verification"
	case usermodel.ActionFinishedCreciFaceVerified:
		iUser := usermodel.NewUser()
		if len(user) > 0 {
			iUser = user[0]
		}
		exist := false
		_, err := us.repo.GetInviteByPhoneNumber(ctx, tx, iUser.GetPhoneNumber())
		if err != nil {
			if status.Code(err) != codes.NotFound {
				return 0, "", 0, err
			}
		} else {
			exist = true
		}

		if exist {
			nextStatus = usermodel.StatusInvitePending
			nextStatusReason = "Awaiting invite verification"
			notification = globalmodel.NotificationRealtorInvitePush
		} else {
			nextStatus = usermodel.StatusActive
			nextStatusReason = "User active"
			notification = globalmodel.NotificationCreciValidated
		}

	case usermodel.ActionFinishedCreciManualVerified:
		nextStatus = usermodel.StatusActive
		nextStatusReason = "User active"
		notification = globalmodel.NotificationCreciValidated
	case usermodel.ActionFinishedCreciImagesUploadedForManualReview:
		nextStatus = usermodel.StatusPendingManual
		nextStatusReason = "Awaiting manual verification by administrator"
	case usermodel.ActionFinishedInviteCreated:
		nextStatus = usermodel.StatusInvitePending
		nextStatusReason = "Awaiting invite verification"
		notification = globalmodel.NotificationRealtorInvitePush
	case usermodel.ActionFinishedInviteAccepted:
		nextStatus = usermodel.StatusActive
		nextStatusReason = "User active"
		notification = globalmodel.NotificationInviteAccepted
	case usermodel.ActionFinishedInviteRejected:
		nextStatus = usermodel.StatusActive
		nextStatusReason = "User active"
		notification = globalmodel.NotificationInviteRejected
	}

	return
}
