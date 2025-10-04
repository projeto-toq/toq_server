package userservices

import (
	"context"
	"database/sql"
	"errors"
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

	if phone := user.GetPhoneNumber(); phone != "" {
		normalizedPhone, normErr := validators.NormalizeToE164(phone)
		if normErr != nil {
			return normErr
		}
		user.SetPhoneNumber(normalizedPhone)
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
		err = utils.ConflictError("Usuário já existe com e-mail, telefone ou CPF informado")
		return
	}

	//verify the password
	if err = validatePassword(user.GetPassword()); err != nil {
		// Padroniza como erro de validação (campo: password)
		return utils.ValidationError("password", "Senha não atende aos requisitos mínimos")
	}

	// Normalize input nationalID before external validation
	if nid := user.GetNationalID(); nid != "" {
		user.SetNationalID(validators.OnlyDigits(nid))
	}

	if role == permissionmodel.RoleSlugAgency {

		cnpj, err1 := us.cnpj.GetCNPJ(ctx, user.GetNationalID())
		if err1 != nil {
			err = us.handleCNPJValidationError(ctx, err1)
			return
		}
		// Ensure digits-only from external provider
		user.SetNationalID(validators.OnlyDigits(cnpj.GetNumeroDeCNPJ()))
		user.SetFullName(cnpj.GetNomeDaPJ())
	} else {
		//validate the userCPF
		cpf, err1 := us.cpf.GetCpf(ctx, user.GetNationalID(), user.GetBornAt())
		if err1 != nil {
			err = us.handleCPFValidationError(ctx, err1)
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
		// Número do endereço é obrigatório
		return utils.ValidationError("number", "Número do endereço é obrigatório")
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

func (us *userService) handleCPFValidationError(ctx context.Context, adapterErr error) error {
	switch {
	case errors.Is(adapterErr, cpfport.ErrInvalidInput):
		slog.Warn("user.validate_user_data.cpf_invalid", "err", adapterErr)
		return utils.ValidationError("national_id", "CPF inválido")
	case errors.Is(adapterErr, cpfport.ErrBirthDateInvalid):
		slog.Warn("user.validate_user_data.cpf_birth_date_invalid", "err", adapterErr)
		return utils.ValidationError("born_at", "Data de nascimento inválida")
	case errors.Is(adapterErr, cpfport.ErrDataMismatch):
		slog.Warn("user.validate_user_data.cpf_birth_date_mismatch", "err", adapterErr)
		return utils.ValidationError("born_at", "Data de nascimento divergente do cadastro da Receita Federal")
	case errors.Is(adapterErr, cpfport.ErrStatusIrregular):
		slog.Warn("user.validate_user_data.cpf_irregular", "err", adapterErr)
		return utils.ValidationError("national_id", "CPF com situação irregular na Receita Federal")
	case errors.Is(adapterErr, cpfport.ErrNotFound):
		slog.Warn("user.validate_user_data.cpf_not_found", "err", adapterErr)
		return utils.ValidationError("national_id", "CPF não encontrado")
	case errors.Is(adapterErr, cpfport.ErrRateLimited):
		slog.Warn("user.validate_user_data.cpf_rate_limited", "err", adapterErr)
		utils.SetSpanError(ctx, adapterErr)
		return utils.TooManyAttemptsError("Limite de consultas ao serviço de CPF atingido")
	case errors.Is(adapterErr, cpfport.ErrInfra):
		slog.Error("user.validate_user_data.cpf_infra_error", "err", adapterErr)
		utils.SetSpanError(ctx, adapterErr)
		return utils.InternalError("Falha ao validar CPF")
	}

	slog.Error("user.validate_user_data.cpf_unhandled_error", "err", adapterErr)
	utils.SetSpanError(ctx, adapterErr)
	return utils.InternalError("Falha ao validar CPF")
}

func (us *userService) handleCNPJValidationError(ctx context.Context, adapterErr error) error {
	switch {
	case errors.Is(adapterErr, cnpjport.ErrInvalid):
		slog.Warn("user.validate_user_data.cnpj_invalid", "err", adapterErr)
		return utils.ValidationError("national_id", "CNPJ inválido")
	case errors.Is(adapterErr, cnpjport.ErrNotFound):
		slog.Warn("user.validate_user_data.cnpj_not_found", "err", adapterErr)
		return utils.ValidationError("national_id", "CNPJ não encontrado")
	case errors.Is(adapterErr, cnpjport.ErrRateLimited):
		slog.Warn("user.validate_user_data.cnpj_rate_limited", "err", adapterErr)
		utils.SetSpanError(ctx, adapterErr)
		return utils.TooManyAttemptsError("Limite de consultas ao serviço de CNPJ atingido")
	case errors.Is(adapterErr, cnpjport.ErrInfra):
		slog.Error("user.validate_user_data.cnpj_infra_error", "err", adapterErr)
		utils.SetSpanError(ctx, adapterErr)
		return utils.InternalError("Falha ao validar CNPJ")
	}

	slog.Error("user.validate_user_data.cnpj_unhandled_error", "err", adapterErr)
	utils.SetSpanError(ctx, adapterErr)
	return utils.InternalError("Falha ao validar CNPJ")
}
