package userservices

import (
	"context"
	"database/sql"
	"time"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateProfileInput represents allowed fields for profile updates.
// Only these fields can be changed by this flow; sensitive fields use dedicated endpoints.
type UpdateProfileInput struct {
	UserID       int64
	NickName     *string
	BornAt       *string // expected format: YYYY-MM-DD
	ZipCode      *string
	Street       *string
	Number       *string
	Complement   *string
	Neighborhood *string
	City         *string
	State        *string // 2-letter state
}

// UpdateProfile updates user profile non-sensitive data.
// It loads the current user, applies provided fields, validates minimal constraints,
// persists changes, and audits the action.
func (us *userService) UpdateProfile(ctx context.Context, in UpdateProfileInput) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("user.update_profile.tracer_error", "error", err)
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Start transaction
	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("user.update_profile.tx_start_error", "error", txErr)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("user.update_profile.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	err = us.updateProfile(ctx, tx, in)
	if err != nil {
		return
	}

	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("user.update_profile.tx_commit_error", "error", commitErr)
		return utils.InternalError("Failed to commit transaction")
	}

	return
}

func (us *userService) updateProfile(
	ctx context.Context,
	tx *sql.Tx,
	in UpdateProfileInput,
) (err error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Carrega o usuário antes de aplicar alterações
	current, err := us.repo.GetUserByID(ctx, tx, in.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.NotFoundError("User")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("user.update_profile.read_user_error", "error", err, "user_id", in.UserID)
		return utils.InternalError("Failed to get user by ID")
	}

	// Aplica somente os campos permitidos se fornecidos
	if in.NickName != nil {
		current.SetNickName(*in.NickName)
	}
	if in.BornAt != nil {
		// valida formato de data YYYY-MM-DD
		bornAt, perr := time.Parse("2006-01-02", *in.BornAt)
		if perr != nil {
			// Validation error with field details
			return utils.ValidationError("born_at", "Invalid date format, expected YYYY-MM-DD")
		}
		current.SetBornAt(bornAt)
	}
	if in.ZipCode != nil {
		current.SetZipCode(*in.ZipCode)
	}
	if in.Street != nil {
		current.SetStreet(*in.Street)
	}
	if in.Number != nil {
		current.SetNumber(*in.Number)
	}
	if in.Complement != nil {
		current.SetComplement(*in.Complement)
	}
	if in.Neighborhood != nil {
		current.SetNeighborhood(*in.Neighborhood)
	}
	if in.City != nil {
		current.SetCity(*in.City)
	}
	if in.State != nil {
		if len(*in.State) != 2 {
			return utils.ValidationError("state", "State must be 2 letters")
		}
		current.SetState(*in.State)
	}

	err = us.repo.UpdateUserByID(ctx, tx, current)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.update_profile.update_user_error", "error", err, "user_id", in.UserID)
		return utils.InternalError("Failed to update user")
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Usuário atualizou o perfil")
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.update_profile.audit_error", "error", err, "user_id", in.UserID)
		return utils.InternalError("Failed to create audit entry")
	}

	return
}
