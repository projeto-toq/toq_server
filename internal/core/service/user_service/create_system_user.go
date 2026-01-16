package userservices

import (
	"context"
	"strings"

	"github.com/google/uuid"

	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	auditservice "github.com/projeto-toq/toq_server/internal/core/service/audit_service"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
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

	nickName := strings.TrimSpace(input.NickName)
	if nickName == "" {
		return SystemUserResult{}, utils.ValidationError("nickName", "Nickname is required")
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

	customZipCode := strings.TrimSpace(input.ZipCode)
	customNumber := strings.TrimSpace(input.Number)
	if customZipCode != "" && customNumber == "" {
		return SystemUserResult{}, utils.ValidationError("number", "Address number is required when zip code is provided")
	}
	if customZipCode == "" && customNumber != "" {
		return SystemUserResult{}, utils.ValidationError("zipCode", "Zip code must be provided when address number is informed")
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
	} else if !errorsIsNoRows(err) {
		utils.SetSpanError(ctx, err)
		logger.Error("admin.users.create.cpf_check_failed", "error", err)
		opErr = utils.InternalError("")
		return SystemUserResult{}, opErr
	}

	passwordSeed := uuid.NewString() + "!Aa1"
	newUser := usermodel.NewUser()
	newUser.SetNickName(nickName)
	newUser.SetEmail(email)
	newUser.SetPhoneNumber(normalizedPhone)
	newUser.SetNationalID(cpfDigits)
	newUser.SetBornAt(input.BornAt)
	newUser.SetPassword(passwordSeed)

	if customZipCode == "" {
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

		newUser.SetZipCode(templateUser.GetZipCode())
		newUser.SetStreet(templateUser.GetStreet())
		newUser.SetNumber(templateUser.GetNumber())
		newUser.SetComplement(templateUser.GetComplement())
		newUser.SetNeighborhood(templateUser.GetNeighborhood())
		newUser.SetCity(templateUser.GetCity())
		newUser.SetState(templateUser.GetState())
	} else {
		newUser.SetZipCode(customZipCode)
		newUser.SetStreet("")
		newUser.SetNumber(customNumber)
		newUser.SetComplement("")
		newUser.SetNeighborhood("")
		newUser.SetCity("")
		newUser.SetState("")
	}

	if err = us.ValidateUserData(ctx, tx, newUser, slug); err != nil {
		opErr = err
		return SystemUserResult{}, opErr
	}

	newUser.SetOptStatus(true)
	newUser.SetNickName(nickName)

	if err = us.repo.CreateUser(ctx, tx, newUser); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("admin.users.create.insert_failed", "error", err)
		opErr = utils.InternalError("")
		return SystemUserResult{}, opErr
	}

	status := globalmodel.StatusActive
	isActive := true
	assignOpts := &AssignRoleOptions{IsActive: &isActive, Status: &status}
	userRole, assignErr := us.AssignRoleToUserWithTx(ctx, tx, newUser.GetID(), role.GetID(), nil, assignOpts)
	if assignErr != nil {
		utils.SetSpanError(ctx, assignErr)
		logger.Error("admin.users.create.assign_role_failed", "user_id", newUser.GetID(), "role_id", role.GetID(), "error", assignErr)
		opErr = utils.InternalError("")
		return SystemUserResult{}, opErr
	}
	newUser.SetActiveRole(userRole)

	if err := us.cloudStorageService.CreateUserFolder(ctx, newUser.GetID()); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("admin.users.create.storage_failed", "user_id", newUser.GetID(), "error", err)
		opErr = utils.InternalError("Failed to prepare user storage")
		return SystemUserResult{}, opErr
	}

	auditRecord := auditservice.BuildRecordFromContext(
		ctx,
		newUser.GetID(),
		auditmodel.AuditTarget{Type: auditmodel.TargetUser, ID: newUser.GetID()},
		auditmodel.OperationCreate,
		map[string]any{
			"role_slug":          role.GetSlug(),
			"origin":             "admin_panel",
			"is_system":          true,
			"phone":              normalizedPhone,
			"cpf":                cpfDigits,
			"zip_code":           newUser.GetZipCode(),
			"has_custom_address": customZipCode != "",
		},
	)
	if auditErr := us.auditService.RecordChange(ctx, tx, auditRecord); auditErr != nil {
		utils.SetSpanError(ctx, auditErr)
		logger.Error("admin.users.create.audit_failed", "user_id", newUser.GetID(), "error", auditErr)
		opErr = utils.InternalError("")
		return SystemUserResult{}, opErr
	}

	// If the new user is a photographer, create their agenda.
	if slug == permissionmodel.RoleSlugPhotographer {
		agendaInput := photosessionservices.EnsureAgendaInput{
			PhotographerID: uint64(newUser.GetID()),
			Timezone:       us.cfg.PhotographerTimezone,
		}
		if agendaErr := us.photoSessionService.EnsurePhotographerAgendaWithTx(ctx, tx, agendaInput); agendaErr != nil {
			utils.SetSpanError(ctx, agendaErr)
			logger.Error("admin.users.create.photographer_agenda_failed", "user_id", newUser.GetID(), "error", agendaErr)
			opErr = utils.InternalError("Failed to create photographer agenda")
			return SystemUserResult{}, opErr
		}
	}

	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("admin.users.create.tx_commit_failed", "error", commitErr)
		return SystemUserResult{}, utils.InternalError("")
	}

	if emailErr := us.sendSystemUserWelcomeEmail(ctx, newUser, slug); emailErr != nil {
		utils.SetSpanError(ctx, emailErr)
		logger.Error("admin.users.create.welcome_email_failed", "user_id", newUser.GetID(), "error", emailErr)
	}

	logger.Info("admin.users.create.success", "user_id", newUser.GetID(), "role_slug", role.GetSlug(), "nick_name", newUser.GetNickName())
	return SystemUserResult{
		UserID:   newUser.GetID(),
		RoleID:   role.GetID(),
		RoleSlug: permissionmodel.RoleSlug(role.GetSlug()),
		Email:    newUser.GetEmail(),
	}, nil
}
