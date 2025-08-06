package userservices

import (
	"context"
	"database/sql"
	"fmt"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) DeleteAgencyOfRealtor(ctx context.Context, realtorID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	err = us.deleteAgencyOfRealtor(ctx, tx, realtorID)
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

func (us *userService) deleteAgencyOfRealtor(ctx context.Context, tx *sql.Tx, realtorID int64) (err error) {

	agency, err := us.repo.GetAgencyOfRealtor(ctx, tx, realtorID)
	if err != nil {
		return
	}

	realtor, err := us.repo.GetUserByID(ctx, tx, realtorID)
	if err != nil {
		return
	}

	_, err = us.repo.DeleteAgencyRealtorRelation(ctx, tx, agency.GetID(), realtor.GetID())
	if err != nil {
		return
	}

	err = us.globalService.SendNotification(ctx, realtor, globalmodel.NotificationAgencyRemovedFromRealtor)
	if err != nil {
		return
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableRealtorAgency,
		fmt.Sprintf("Apagado o relacionamento com a Imobiliária %s", agency.GetNickName()))
	if err != nil {
		return
	}

	return
}
