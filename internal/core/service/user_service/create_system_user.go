package userservices

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	permissionservices "github.com/projeto-toq/toq_server/internal/core/service/permission_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
	validators "github.com/projeto-toq/toq_server/internal/core/utils/validators"
)

// CreateSystemUser cria um usuário de sistema replicando endereço do usuário raiz (ID=1) e atribuindo role sistêmico.
func (us *userService) CreateSystemUser(ctx context.Context, input CreateSystemUserInput) (SystemUserResult, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return SystemUserResult{}, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

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

	cpfDigits := validators.OnlyDigits(input.CPF)
	if cpfDigits == "" {
		return SystemUserResult{}, utils.ValidationError("cpf", "CPF is required")
	}

	if input.BornAt.IsZero() {
		return SystemUserResult{}, utils.ValidationError("bornAt", "Birth date is required")
	}

	slug := permissionmodel.RoleSlug(strings.TrimSpace(input.RoleSlug.String()))
	if slug == "" {
		return SystemUserResult{}, utils.ValidationError("roleSlug", "Role slug is required")
	}
	if !slug.IsValid() {
		return SystemUserResult{}, utils.ValidationError("roleSlug", "Invalid role slug")
	}

	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("admin.users.create.tx_start_failed", "error", txErr)
		return SystemUserResult{}, utils.InternalError("")
	}

	var opErr error
	defer func() {
		if opErr != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("admin.users.create.tx_rollback_failed", "error", rbErr)
			}
		}
	}()

	role, roleErr := us.permissionService.GetRoleBySlugWithTx(ctx, tx, slug)
	if roleErr != nil {
		utils.SetSpanError(ctx, roleErr)
		logger.Error("admin.users.create.get_role_failed", "slug", slug, "error", roleErr)
		opErr = utils.InternalError("")
		return SystemUserResult{}, opErr
	}
	if role == nil {
		opErr = utils.NotFoundError("role")
		return SystemUserResult{}, opErr
	}
	if !role.GetIsSystemRole() {
		opErr = derrors.ErrRoleNotSystem
		return SystemUserResult{}, opErr
	}

	emailExists, emailCheckErr := us.repo.ExistsEmailForAnotherUser(ctx, tx, email, 0)
	if emailCheckErr != nil {
		utils.SetSpanError(ctx, emailCheckErr)
		logger.Error("admin.users.create.email_check_failed", "error", emailCheckErr)
		opErr = utils.InternalError("")
		return SystemUserResult{}, opErr
	}
	if emailExists {
		opErr = utils.ConflictError("Email already in use")
		return SystemUserResult{}, opErr
	}

	phoneExists, phoneCheckErr := us.repo.ExistsPhoneForAnotherUser(ctx, tx, normalizedPhone, 0)
	if phoneCheckErr != nil {
		utils.SetSpanError(ctx, phoneCheckErr)
		logger.Error("admin.users.create.phone_check_failed", "error", phoneCheckErr)
		opErr = utils.InternalError("")
		return SystemUserResult{}, opErr
	}
	if phoneExists {
		opErr = utils.ConflictError("Phone already in use")
		return SystemUserResult{}, opErr
	}

	if _, err := us.repo.GetUserByNationalID(ctx, tx, cpfDigits); err == nil {
		opErr = utils.ConflictError("CPF already in use")
		return SystemUserResult{}, opErr
	} else if err != nil && !errorsIsNoRows(err) {
		utils.SetSpanError(ctx, err)
		logger.Error("admin.users.create.cpf_check_failed", "error", err)
		opErr = utils.InternalError("")
		return SystemUserResult{}, opErr
	}

	templateUser, templateErr := us.repo.GetUserByID(ctx, tx, systemUserTemplateID)
	if templateErr != nil {
		utils.SetSpanError(ctx, templateErr)
		logger.Error("admin.users.create.template_fetch_failed", "template_id", systemUserTemplateID, "error", templateErr)
		if errorsIsNoRows(templateErr) {
			opErr = utils.InternalError("System template user not found")
		} else {
			opErr = utils.InternalError("")
		}
		return SystemUserResult{}, opErr
	}

	now := time.Now().UTC()
	newUser := usermodel.NewUser()
	newUser.SetFullName(fullName)
	newUser.SetNickName(firstToken(fullName))
	newUser.SetNationalID(cpfDigits)
	newUser.SetBornAt(input.BornAt)
	newUser.SetPhoneNumber(normalizedPhone)
	newUser.SetEmail(email)
	newUser.SetOptStatus(false)
	newUser.SetDeleted(false)
	newUser.SetLastActivityAt(now)
	newUser.SetLastSignInAttempt(time.Time{})

	newUser.SetZipCode(templateUser.GetZipCode())
	newUser.SetStreet(templateUser.GetStreet())
	newUser.SetNumber(templateUser.GetNumber())
	newUser.SetComplement(templateUser.GetComplement())
	newUser.SetNeighborhood(templateUser.GetNeighborhood())
	newUser.SetCity(templateUser.GetCity())
	newUser.SetState(templateUser.GetState())

	passwordSeed := uuid.NewString() + "!Aa1"
	newUser.SetPassword(us.encryptPassword(passwordSeed))

	if err = us.repo.CreateUser(ctx, tx, newUser); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("admin.users.create.insert_failed", "error", err)
		opErr = utils.InternalError("")
		return SystemUserResult{}, opErr
	}

	status := permissionmodel.StatusActive
	isActive := true
	assignOpts := &permissionservices.AssignRoleOptions{IsActive: &isActive, Status: &status}
	userRole, assignErr := us.permissionService.AssignRoleToUserWithTx(ctx, tx, newUser.GetID(), role.GetID(), nil, assignOpts)
	if assignErr != nil {
		utils.SetSpanError(ctx, assignErr)
		logger.Error("admin.users.create.assign_role_failed", "user_id", newUser.GetID(), "role_id", role.GetID(), "error", assignErr)
		opErr = utils.InternalError("")
		return SystemUserResult{}, opErr
	}
	newUser.SetActiveRole(userRole)

	if auditErr := us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Criado usuário do sistema via painel admin", newUser.GetID()); auditErr != nil {
		utils.SetSpanError(ctx, auditErr)
		logger.Error("admin.users.create.audit_failed", "user_id", newUser.GetID(), "error", auditErr)
		opErr = utils.InternalError("")
		return SystemUserResult{}, opErr
	}

	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("admin.users.create.tx_commit_failed", "error", commitErr)
		return SystemUserResult{}, utils.InternalError("")
	}

	logger.Info("admin.users.create.success", "user_id", newUser.GetID(), "role_slug", role.GetSlug())
	return SystemUserResult{
		UserID:   newUser.GetID(),
		RoleID:   role.GetID(),
		RoleSlug: permissionmodel.RoleSlug(role.GetSlug()),
		Email:    newUser.GetEmail(),
	}, nil
}
