package userservices

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	derrors "github.com/giulio-alfieri/toq_server/internal/core/derrors"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// decideNextStatusAfterContactChange decides the next status after a successful email/phone change.
// It is a pure function: given the role and pending flags, it returns the target status.
func decideNextStatusAfterContactChange(role permissionmodel.RoleSlug, emailPending, phonePending bool, current permissionmodel.UserRoleStatus) permissionmodel.UserRoleStatus {
	// Primeiro, se ainda há pendência de algum fator de contato, refletir essa pendência no status.
	if emailPending {
		return permissionmodel.StatusPendingEmail
	}
	if phonePending {
		return permissionmodel.StatusPendingPhone
	}
	// Caso não haja mais pendências, decidir por role.
	switch role {
	case permissionmodel.RoleSlugOwner:
		return permissionmodel.StatusActive
	case permissionmodel.RoleSlugRealtor:
		return permissionmodel.StatusPendingCreci
	case permissionmodel.RoleSlugAgency:
		return permissionmodel.StatusPendingCnpj
	default:
		// Papel desconhecido: manter status atual por segurança.
		return current
	}
}

// applyStatusTransitionAfterContactChange loads user and validations, then applies the target status
// based on pending contact factors and active role. Runs inside the provided transaction.
func (us *userService) applyStatusTransitionAfterContactChange(ctx context.Context, tx *sql.Tx, emailJustConfirmed bool) (permissionmodel.UserRoleStatus, bool, error) {
	// Carregar usuário e papel ativo
	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		// domínio: auth
		return 0, false, derrors.Auth("Authentication required")
	}
	// Resolver papel ativo via permission service (robusto e desacoplado de carregamento do usuário)
	active, derr := us.permissionService.GetActiveUserRoleWithTx(ctx, tx, userID)
	if derr != nil {
		// pode ser domínio (KindError) ou erro puro de infra – deixar o caller mapear/logar
		return 0, false, derr
	}
	if active == nil || active.GetRole() == nil {
		return 0, false, derrors.ErrUserActiveRoleMissing
	}
	roleSlug := permissionmodel.RoleSlug(active.GetRole().GetSlug())
	from := active.GetStatus()

	// Ler validações para checar pendências do outro fator
	validations, err := us.repo.GetUserValidations(ctx, tx, userID)
	var emailPending, phonePending bool
	if err != nil {
		// Ausência de validações após confirmação significa "sem pendências" e não é erro
		if errors.Is(err, sql.ErrNoRows) {
			emailPending, phonePending = false, false
		} else {
			// Outros erros são infraestrutura
			slog.Error("user.status_transition.stage_error", "stage", "get_validations", "err", err)
			return 0, false, err
		}
	} else {
		emailPending = validations.GetEmailCode() != ""
		phonePending = validations.GetPhoneCode() != ""
	}
	// Ajustar o fator recém confirmado
	if emailJustConfirmed {
		emailPending = false
	} else {
		phonePending = false
	}

	to := decideNextStatusAfterContactChange(roleSlug, emailPending, phonePending, from)
	if to == from {
		return from, false, nil
	}

	// Persistir alteração e auditar
	if err := us.repo.UpdateUserRoleStatus(ctx, tx, userID, roleSlug, to); err != nil {
		// Se não houve role ativo correspondente, mapear para erro de domínio e não logar como infra
		if errors.Is(err, sql.ErrNoRows) {
			return 0, false, derrors.ErrUserActiveRoleMissing
		}
		slog.Error("user.status_transition.stage_error", "stage", "update_role_status", "err", err)
		return 0, false, err // infra
	}
	if err := us.globalService.CreateAudit(ctx, tx, globalmodel.TableUserRoles, "Atualização de status após alteração de contato"); err != nil {
		slog.Error("user.status_transition.stage_error", "stage", "audit", "err", err)
		return 0, false, err // infra
	}

	return to, true, nil
}
