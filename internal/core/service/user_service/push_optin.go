package userservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) PushOptIn(ctx context.Context, userID int64, deviceToken string) (err error) {
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

	err = us.pushOptIn(ctx, tx, userID, deviceToken)
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

func (us *userService) pushOptIn(ctx context.Context, tx *sql.Tx, userID int64, deviceToken string) (err error) {
	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		return
	}
	//TODO: Remove this hardcoded deviceToken
	s := "dDWfs2iRThyJvzd_dSvyah:APA91bGp1GdU1zNsTzpaNb9gJpPdPTOVvJFpL2vpT52E7wemRocGtCe8HN5rpxk_Ys5NH4qo__7CD4_TZ0ahbTk2CyRaj36gCwlV9IANjFFtiQpQEvbSenw"
	user.SetDeviceToken(s)
	_ = deviceToken

	//TODO: Uncomment this line to use the deviceToken from the request
	// user.SetDeviceToken(deviceToken)

	err = us.repo.UpdateUserByID(ctx, tx, user)
	if err != nil {
		return
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Usuário aceitou receber notificações")
	if err != nil {
		return
	}

	return
}
