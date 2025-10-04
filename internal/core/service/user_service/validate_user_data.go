package userservices

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	cnpjport "github.com/giulio-alfieri/toq_server/internal/core/port/right/cnpj"
	cpfport "github.com/giulio-alfieri/toq_server/internal/core/port/right/cpf"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	validators "github.com/giulio-alfieri/toq_server/internal/core/utils/validators"
)

func (us *userService) ValidateUserData(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface, role permissionmodel.RoleSlug) (err error) {

	now := time.Now().UTC()
	prefix := dataPrefixForRole(role)
	field := func(name string) string {
		return composeField(prefix, name)
	}

	if phone := user.GetPhoneNumber(); phone != "" {
		normalizedPhone, normErr := validators.NormalizeToE164(phone)
		if normErr != nil {
			return utils.ValidationError(field("phoneNumber"), "Invalid phone number format.")
		}
		user.SetPhoneNumber(normalizedPhone)
	}

	if email := user.GetEmail(); email != "" {
		if err := validators.ValidateEmail(email); err != nil {
			return utils.ValidationError(field("email"), "Invalid email format.")
		}
	}

	//verify if user already exists
	exist, err := us.repo.VerifyUserDuplicity(ctx, tx, user)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.validate_user_data.duplicity_check_error", "error", err)
		return
	}

	if exist {
		// Usuário existente (email/telefone/CPF)
		err = utils.ConflictError("User already exists with provided email, phone or national ID.")
		return
	}

	//verify the password
	if err = validatePassword(field("password"), user.GetPassword()); err != nil {
		return err
	}

	// Normalize input nationalID before external validation
	if nid := user.GetNationalID(); nid != "" {
		user.SetNationalID(validators.OnlyDigits(nid))
	}

	if role == permissionmodel.RoleSlugAgency {

		cnpj, err1 := us.cnpj.GetCNPJ(ctx, user.GetNationalID())
		if err1 != nil {
			err = us.handleCNPJValidationError(ctx, prefix, err1)
			return
		}
		// Ensure digits-only from external provider
		user.SetNationalID(validators.OnlyDigits(cnpj.GetNumeroDeCNPJ()))
		user.SetFullName(cnpj.GetNomeDaPJ())
	} else {
		//validate the userCPF
		cpf, err1 := us.cpf.GetCpf(ctx, user.GetNationalID(), user.GetBornAt())
		if err1 != nil {
			err = us.handleCPFValidationError(ctx, prefix, err1)
			return
		}
		// Ensure digits-only from external provider
		user.SetNationalID(validators.OnlyDigits(cpf.GetNumeroDeCpf()))
		user.SetFullName(cpf.GetNomeDaPf())
	}

	//validate the user zipcode
	cep, err := us.globalService.GetCEP(ctx, user.GetZipCode())
	if err != nil {
		return err
	}

	//validate the address number
	if user.GetNumber() == "" {
		// Address number is required
		return utils.ValidationError(field("number"), "Address number is required.")
	}

	user.SetStreet(cep.GetStreet())
	// user.SetComplement(cep.GetComplement())
	user.SetNeighborhood(cep.GetNeighborhood())
	user.SetCity(cep.GetCity())
	user.SetState(cep.GetState())

	user.SetPassword(us.encryptPassword(user.GetPassword()))
	user.SetLastActivityAt(now)
	user.SetDeleted(false)

	return
}

func (us *userService) handleCPFValidationError(ctx context.Context, prefix string, adapterErr error) error {
	field := func(name string) string {
		return composeField(prefix, name)
	}
	switch {
	case errors.Is(adapterErr, cpfport.ErrInvalidInput):
		slog.Warn("user.validate_user_data.cpf_invalid", "err", adapterErr)
		return utils.ValidationError(field("nationalID"), "Invalid national ID.")
	case errors.Is(adapterErr, cpfport.ErrBirthDateInvalid):
		slog.Warn("user.validate_user_data.cpf_birth_date_invalid", "err", adapterErr)
		return utils.ValidationError(field("bornAt"), "Invalid birth date.")
	case errors.Is(adapterErr, cpfport.ErrDataMismatch):
		slog.Warn("user.validate_user_data.cpf_birth_date_mismatch", "err", adapterErr)
		return utils.ValidationError(field("bornAt"), "Birth date does not match government records.")
	case errors.Is(adapterErr, cpfport.ErrStatusIrregular):
		slog.Warn("user.validate_user_data.cpf_irregular", "err", adapterErr)
		return utils.ValidationError(field("nationalID"), "National ID has an irregular status.")
	case errors.Is(adapterErr, cpfport.ErrNotFound):
		slog.Warn("user.validate_user_data.cpf_not_found", "err", adapterErr)
		return utils.ValidationError(field("nationalID"), "National ID not found.")
	case errors.Is(adapterErr, cpfport.ErrRateLimited):
		slog.Warn("user.validate_user_data.cpf_rate_limited", "err", adapterErr)
		utils.SetSpanError(ctx, adapterErr)
		return utils.TooManyAttemptsError("National ID lookup rate limit exceeded.")
	case errors.Is(adapterErr, cpfport.ErrInfra):
		slog.Error("user.validate_user_data.cpf_infra_error", "err", adapterErr)
		utils.SetSpanError(ctx, adapterErr)
		return utils.InternalError("Failed to validate national ID.")
	}

	slog.Error("user.validate_user_data.cpf_unhandled_error", "err", adapterErr)
	utils.SetSpanError(ctx, adapterErr)
	return utils.InternalError("Failed to validate national ID.")
}

func (us *userService) handleCNPJValidationError(ctx context.Context, prefix string, adapterErr error) error {
	field := func(name string) string {
		return composeField(prefix, name)
	}
	switch {
	case errors.Is(adapterErr, cnpjport.ErrInvalid):
		slog.Warn("user.validate_user_data.cnpj_invalid", "err", adapterErr)
		return utils.ValidationError(field("nationalID"), "Invalid company national ID.")
	case errors.Is(adapterErr, cnpjport.ErrNotFound):
		slog.Warn("user.validate_user_data.cnpj_not_found", "err", adapterErr)
		return utils.ValidationError(field("nationalID"), "Company national ID not found.")
	case errors.Is(adapterErr, cnpjport.ErrRateLimited):
		slog.Warn("user.validate_user_data.cnpj_rate_limited", "err", adapterErr)
		utils.SetSpanError(ctx, adapterErr)
		return utils.TooManyAttemptsError("Company national ID lookup rate limit exceeded.")
	case errors.Is(adapterErr, cnpjport.ErrInfra):
		slog.Error("user.validate_user_data.cnpj_infra_error", "err", adapterErr)
		utils.SetSpanError(ctx, adapterErr)
		return utils.InternalError("Failed to validate company national ID.")
	}

	slog.Error("user.validate_user_data.cnpj_unhandled_error", "err", adapterErr)
	utils.SetSpanError(ctx, adapterErr)
	return utils.InternalError("Failed to validate company national ID.")
}

func dataPrefixForRole(role permissionmodel.RoleSlug) string {
	switch role {
	case permissionmodel.RoleSlugOwner:
		return "owner"
	case permissionmodel.RoleSlugRealtor:
		return "realtor"
	case permissionmodel.RoleSlugAgency:
		return "agency"
	default:
		return "user"
	}
}

func composeField(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return fmt.Sprintf("%s.%s", prefix, name)
}
