package userservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/validators"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (us *userService) ValidateCreciData(ctx context.Context, realtors []usermodel.UserInterface) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Start a database transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	err = us.validateCreciData(ctx, tx, realtors)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

}

func (us *userService) validateCreciData(
	ctx context.Context,
	tx *sql.Tx,
	realtors []usermodel.UserInterface,
) (err error) {
	auditMessage := ""
	for _, realtor := range realtors {

		// Call the OCR
		creci, err1 := us.creci.ValidateCreciNumber(ctx, realtor)
		if err1 != nil {
			if status.Code(err1) == codes.InvalidArgument {
				status, reason, notification, err2 := us.updateUserStatus(ctx, tx, realtor.GetActiveRole().GetRole(), usermodel.ActionFinishedCreciStateUnsupported)
				if err2 != nil {
					return err2
				}
				realtor.GetActiveRole().SetStatus(status)
				realtor.GetActiveRole().SetStatusReason(reason)
				err2 = us.globalService.SendNotification(ctx, realtor, notification)
				if err2 != nil {
					return err2
				}
				auditMessage = "Creci de estado não suportado"
			} else {
				return err1
			}
		}
		//TODO: incluir validação de campo nulo antes de chamar as comparações
		if !validators.ValidateCreciEquality(creci.GetCreciNumber(), realtor.GetCreciNumber()) {
			status, reason, notification, err2 := us.updateUserStatus(ctx, tx, realtor.GetActiveRole().GetRole(), usermodel.ActionFinishedCreciNumberDoesntMatch)
			if err2 != nil {
				return err2
			}
			realtor.GetActiveRole().SetStatus(status)
			realtor.GetActiveRole().SetStatusReason(reason)
			err1 = us.globalService.SendNotification(ctx, realtor, notification)
			if err1 != nil {
				return err1
			}
			auditMessage = "Numero do creci não confere"
		} else if realtor.GetCreciState() != creci.GetCreciState() {
			status, reason, notification, err2 := us.updateUserStatus(ctx, tx, realtor.GetActiveRole().GetRole(), usermodel.ActionFinishedCreciStateDoesntMatch)
			if err2 != nil {
				return err2
			}
			realtor.GetActiveRole().SetStatus(status)
			realtor.GetActiveRole().SetStatusReason(reason)
			err1 = us.globalService.SendNotification(ctx, realtor, notification)
			if err1 != nil {
				return err1
			}
			auditMessage = "Estado do Creci não confere"
		} else {
			status, reason, _, err2 := us.updateUserStatus(ctx, tx, realtor.GetActiveRole().GetRole(), usermodel.ActionFinishedCreciVerified)
			if err2 != nil {
				return err2
			}
			realtor.GetActiveRole().SetStatus(status)
			realtor.GetActiveRole().SetStatusReason(reason)
			auditMessage = "Creci do corretor validado"
		}
		err = us.repo.UpdateUserByID(ctx, tx, realtor)
		if err != nil {
			return
		}
		err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, auditMessage)
		if err != nil {
			return
		}
	}

	return
}
