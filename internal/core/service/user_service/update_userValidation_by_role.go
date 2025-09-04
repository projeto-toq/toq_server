package userservices

import (
	"context"
	"database/sql"
	"log/slog"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// UpdateUserValidationByUserRole is deprecated due to unsafe dereferencing of ActiveRole/Role.
// Use UpdateUserValidationByRole instead. Mantido temporariamente para compatibilidade interna.
func (us *userService) UpdateUserValidationByUserRole(ctx context.Context, tx *sql.Tx, user *usermodel.UserInterface, userValidation usermodel.ValidationInterface) (bool, error) {
	// Redireciona para a versão segura utilizando o userID a partir da validação.
	return us.UpdateUserValidationByRole(ctx, tx, userValidation.GetUserID(), userValidation)
}

// UpdateUserValidationByRole safely updates user status after contact validation using the active role
// resolved via permission service, avoiding nil dereferences.
// Returns whether tokens should be generated after the update.
func (us *userService) UpdateUserValidationByRole(ctx context.Context, tx *sql.Tx, userID int64, userValidation usermodel.ValidationInterface) (generateTokens bool, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return false, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	generateTokens = false

	// Obter role ativa de forma segura via permission service (dentro da mesma transação)
	activeRole, aerr := us.permissionService.GetActiveUserRoleWithTx(ctx, tx, userID)
	if aerr != nil {
		if derr, ok := aerr.(utils.DomainError); ok {
			return false, utils.WrapDomainErrorWithSource(derr)
		}
		return false, utils.InternalError("Failed to get active user role")
	}
	if activeRole == nil || activeRole.GetRole() == nil {
		slog.Warn("user.validation.active_role_missing", "user_id", userID)
		return false, utils.ErrUserActiveRoleMissing
	}

	roleSlug := utils.GetUserRoleSlugFromUserRole(activeRole)

	// Decidir qual ação de status executar conforme validações pendentes
	switch {
	case userValidation.GetEmailCode() == "" && userValidation.GetPhoneCode() == "":
		// ambos validados → perfil pronto
		if _, _, _, err = us.updateUserStatus(ctx, tx, roleSlug, usermodel.ActionFinishedEmailVerified); err != nil {
			if derr, ok := err.(utils.DomainError); ok {
				return false, utils.WrapDomainErrorWithSource(derr)
			}
			return false, utils.InternalError("Failed to update user status")
		}
		generateTokens = true

	case userValidation.GetEmailCode() == "" && userValidation.GetPhoneCode() != "":
		// apenas telefone validado → pendente email
		if _, _, _, err = us.updateUserStatus(ctx, tx, roleSlug, usermodel.ActionFinishedPhoneVerified); err != nil {
			if derr, ok := err.(utils.DomainError); ok {
				return false, utils.WrapDomainErrorWithSource(derr)
			}
			return false, utils.InternalError("Failed to update user status")
		}
		generateTokens = false

	default:
		// Nenhuma transição necessária neste momento
		generateTokens = false
	}

	return generateTokens, nil
}
