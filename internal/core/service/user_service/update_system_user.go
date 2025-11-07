package userservices

import (
	"context"
	"strings"

	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
	validators "github.com/projeto-toq/toq_server/internal/core/utils/validators"
)

// UpdateSystemUser atualiza dados sensíveis de um usuário de sistema.
func (us *userService) UpdateSystemUser(ctx context.Context, input UpdateSystemUserInput) (SystemUserResult, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return SystemUserResult{}, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.UserID <= 0 {
		return SystemUserResult{}, utils.ValidationError("userId", "User id must be positive")
	}

	fullName := strings.TrimSpace(input.FullName)
	if fullName == "" {
		return SystemUserResult{}, utils.ValidationError("fullName", "Full name is required")
	}

	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" {
		return SystemUserResult{}, utils.ValidationError("email", "Email is required")
	}
	if err := validators.ValidateEmail(email); err != nil {
		return SystemUserResult{}, err
	}

	normalizedPhone, phoneErr := validators.NormalizeToE164(input.PhoneNumber)
	if phoneErr != nil {
		return SystemUserResult{}, phoneErr
	}

	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("admin.users.update.tx_start_failed", "error", txErr)
		return SystemUserResult{}, utils.InternalError("")
	}

	var opErr error
	defer func() {
		if opErr != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("admin.users.update.tx_rollback_failed", "error", rbErr)
			}
		}
	}()

	existing, userErr := us.repo.GetUserByID(ctx, tx, input.UserID)
	if userErr != nil {
		utils.SetSpanError(ctx, userErr)
		logger.Error("admin.users.update.get_user_failed", "user_id", input.UserID, "error", userErr)
		if errorsIsNoRows(userErr) {
			opErr = utils.NotFoundError("user")
		} else {
			opErr = utils.InternalError("")
		}
		return SystemUserResult{}, opErr
	}

	if existing.IsDeleted() {
		opErr = derrors.ErrUserAlreadyDeleted
		return SystemUserResult{}, opErr
	}

	activeRole, arErr := us.GetActiveUserRoleWithTx(ctx, tx, input.UserID)
	if arErr != nil {
		utils.SetSpanError(ctx, arErr)
		logger.Error("admin.users.update.get_role_failed", "user_id", input.UserID, "error", arErr)
		opErr = utils.InternalError("")
		return SystemUserResult{}, opErr
	}
	if activeRole == nil || activeRole.GetRole() == nil || !activeRole.GetRole().GetIsSystemRole() {
		opErr = derrors.ErrSystemUserRoleMismatch
		return SystemUserResult{}, opErr
	}

	emailExists, emailCheckErr := us.repo.ExistsEmailForAnotherUser(ctx, tx, email, input.UserID)
	if emailCheckErr != nil {
		utils.SetSpanError(ctx, emailCheckErr)
		logger.Error("admin.users.update.email_check_failed", "user_id", input.UserID, "error", emailCheckErr)
		opErr = utils.InternalError("")
		return SystemUserResult{}, opErr
	}
	if emailExists {
		opErr = utils.ConflictError("Email already in use")
		return SystemUserResult{}, opErr
	}

	phoneExists, phoneCheckErr := us.repo.ExistsPhoneForAnotherUser(ctx, tx, normalizedPhone, input.UserID)
	if phoneCheckErr != nil {
		utils.SetSpanError(ctx, phoneCheckErr)
		logger.Error("admin.users.update.phone_check_failed", "user_id", input.UserID, "error", phoneCheckErr)
		opErr = utils.InternalError("")
		return SystemUserResult{}, opErr
	}
	if phoneExists {
		opErr = utils.ConflictError("Phone already in use")
		return SystemUserResult{}, opErr
	}

	existing.SetFullName(fullName)
	existing.SetNickName(firstToken(fullName))
	existing.SetEmail(email)
	existing.SetPhoneNumber(normalizedPhone)

	if updateErr := us.repo.UpdateUserByID(ctx, tx, existing); updateErr != nil {
		utils.SetSpanError(ctx, updateErr)
		logger.Error("admin.users.update.repo_failed", "user_id", input.UserID, "error", updateErr)
		opErr = utils.InternalError("")
		return SystemUserResult{}, opErr
	}

	if auditErr := us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Atualizado usuário do sistema via painel admin", existing.GetID()); auditErr != nil {
		utils.SetSpanError(ctx, auditErr)
		logger.Error("admin.users.update.audit_failed", "user_id", input.UserID, "error", auditErr)
		opErr = utils.InternalError("")
		return SystemUserResult{}, opErr
	}

	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("admin.users.update.tx_commit_failed", "user_id", input.UserID, "error", commitErr)
		return SystemUserResult{}, utils.InternalError("")
	}

	logger.Info("admin.users.update.success", "user_id", existing.GetID())
	return SystemUserResult{
		UserID:   existing.GetID(),
		RoleID:   activeRole.GetRoleID(),
		RoleSlug: permissionmodel.RoleSlug(activeRole.GetRole().GetSlug()),
		Email:    existing.GetEmail(),
	}, nil
}
