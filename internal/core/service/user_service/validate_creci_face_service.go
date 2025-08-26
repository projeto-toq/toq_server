package userservices

import (
	"context"
	"database/sql"
	"log/slog"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) ValidateCreciFace(ctx context.Context, realtors []usermodel.UserInterface) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	err = us.validateCreciFace(ctx, tx, realtors)
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

}

func (us *userService) validateCreciFace(
	ctx context.Context,
	tx *sql.Tx,

	realtors []usermodel.UserInterface,
) (err error) {
	// Verificar se CRECI adapter está disponível
	if us.creci == nil {
		slog.Warn("CRECI adapter não disponível, pulando validação de face")
		return nil
	}

	auditMessage := ""
	for _, realtor := range realtors {
		valid, err1 := us.creci.ValidateFaceMatch(ctx, realtor)
		if err1 != nil {
			status, reason, notification, err2 := us.updateUserStatus(ctx, tx, realtor.GetActiveRole().GetRole(), usermodel.ActionFinishedBadSelfieImage)
			if err2 != nil {
				return err2
			}
			realtor.GetActiveRole().SetStatus(status)
			realtor.GetActiveRole().SetStatusReason(reason)
			err1 := us.globalService.SendNotification(ctx, realtor, notification)
			if err1 != nil {
				return err1
			}
			err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "A selfie está com baixa qualidade")
			if err != nil {
				return
			}
			return
		}

		if valid {
			status, reason, notification, err2 := us.updateUserStatus(ctx, tx, realtor.GetActiveRole().GetRole(), usermodel.ActionFinishedCreciFaceVerified, realtor)
			if err2 != nil {
				return err2
			}
			realtor.GetActiveRole().SetStatus(status)
			realtor.GetActiveRole().SetStatusReason(reason)
			err1 := us.globalService.SendNotification(ctx, realtor, notification)
			if err1 != nil {
				return err1
			}
			auditMessage = "Selfie do corretor validado"

		} else {
			status, reason, notification, err2 := us.updateUserStatus(ctx, tx, realtor.GetActiveRole().GetRole(), usermodel.ActionFinishedSelfieDoesntMatch)
			if err2 != nil {
				return err2
			}
			realtor.GetActiveRole().SetStatus(status)
			realtor.GetActiveRole().SetStatusReason(reason)

			err1 := us.globalService.SendNotification(ctx, realtor, notification)
			if err1 != nil {
				return err1
			}
			auditMessage = "Selfie do corretor rejeitada"
		}

		err2 := us.repo.UpdateUserByID(ctx, tx, realtor)
		if err2 != nil {
			return err2
		}
		err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, auditMessage)
		if err != nil {
			return
		}

	}

	return
}
