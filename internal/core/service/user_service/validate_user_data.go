package userservices

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	validators "github.com/giulio-alfieri/toq_server/internal/core/utils/validators"
)

func (us *userService) ValidateUserData(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface, role permissionmodel.RoleSlug) (err error) {

	now := time.Now().UTC()

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

		cnpj, err1 := us.cnpj.GetCNPJ(ctx, user.GetNationalID()) // external validation
		if err1 != nil {
			// Propaga erro do adaptador (serviço externo ou dado inválido)
			utils.SetSpanError(ctx, err1)
			slog.Error("user.validate_user_data.cnpj_error", "error", err1)
			err = err1
			return
		}
		// Ensure digits-only from external provider
		user.SetNationalID(validators.OnlyDigits(cnpj.GetNumeroDeCNPJ()))
		user.SetFullName(cnpj.GetNomeDaPJ())
	} else {
		//validate the userCPF
		cpf, err1 := us.cpf.GetCpf(ctx, user.GetNationalID(), user.GetBornAt())
		if err1 != nil {
			// Propaga erro do adaptador (serviço externo ou dado inválido)
			utils.SetSpanError(ctx, err1)
			slog.Error("user.validate_user_data.cpf_error", "error", err1)
			err = err1
			return
		}
		// Ensure digits-only from external provider
		user.SetNationalID(validators.OnlyDigits(cpf.GetNumeroDeCpf()))
		user.SetFullName(cpf.GetNomeDaPf())
	}

	//validate the user zipcode
	cep, err := us.globalService.GetCEP(ctx, user.GetZipCode())
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.validate_user_data.cep_error", "error", err)
		return utils.InternalError("Failed to validate zipcode")
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
