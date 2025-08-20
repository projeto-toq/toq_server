package userservices

import (
	"context"
	"database/sql"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (us *userService) VerifyCreciImages(ctx context.Context, realtorID int64) (err error) {

	//create a new transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}
	err = us.verifyCreciImages(ctx, tx, realtorID)
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

func (us *userService) verifyCreciImages(ctx context.Context, tx *sql.Tx, realtorID int64) (err error) {

	//verify the realtor is awaiting images
	realtor, err := us.repo.GetUserByID(ctx, tx, realtorID)
	if err != nil {
		return
	}

	if !(realtor.GetActiveRole().GetStatus() == usermodel.StatusPendingImages ||
		realtor.GetActiveRole().GetStatus() == usermodel.StatusRejectByOCR ||
		realtor.GetActiveRole().GetStatus() == usermodel.StatusRejectByFace ||
		realtor.GetActiveRole().GetStatus() != usermodel.StatusPendingManual) {
		err = status.Error(codes.FailedPrecondition, "user is not awaiting creci images")
		return
	}

	status, reason, _, err := us.updateUserStatus(ctx, tx, realtor.GetActiveRole().GetRole(), usermodel.ActionFinishedCreciImagesUploadedForManualReview)
	if err != nil {
		return
	}

	realtor.GetActiveRole().SetStatus(status)
	realtor.GetActiveRole().SetStatusReason(reason)

	//update the status of the user to pendingOCR
	err = us.repo.UpdateUserByID(ctx, tx, realtor)
	if err != nil {
		return
	}
	return
}
