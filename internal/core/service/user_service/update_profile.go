package userservices

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
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
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	// Start transaction
	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		slog.Error("user.update_profile.tx_start_error", "err", txErr)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("user.update_profile.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	err = us.updateProfile(ctx, tx, in)
	if err != nil {
		return
	}

	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		slog.Error("user.update_profile.tx_commit_error", "err", commitErr)
		return utils.InternalError("Failed to commit transaction")
	}

	return
}

func (us *userService) updateProfile(
	ctx context.Context,
	tx *sql.Tx,
	in UpdateProfileInput,
) (err error) {
	// Carrega o usuário antes de aplicar alterações
	current, err := us.repo.GetUserByID(ctx, tx, in.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.NotFoundError("User")
		}
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
		return utils.InternalError("Failed to update user")
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Usuário atualizou o perfil")
	if err != nil {
		return utils.InternalError("Failed to create audit entry")
	}

	return
}
