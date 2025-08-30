package userservices

import (
	"context"
	"database/sql"
	"time"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) ValidateUserData(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface, role permissionmodel.RoleSlug) (err error) {

	now := time.Now().UTC()

	//verify if user already exists
	exist, err := us.repo.VerifyUserDuplicity(ctx, tx, user)
	if err != nil {
		return
	}

	if exist {
		err = utils.ErrInternalServer
		return
	}

	//verify the password
	err = validatePassword(user.GetPassword())
	if err != nil {
		return
	}

	if role == permissionmodel.RoleSlugAgency {

		cnpj, err1 := us.cnpj.GetCNPJ(ctx, user.GetNationalID()) // Validation via global service integration planned
		if err1 != nil {
			err = err1
			return
		}
		user.SetNationalID(cnpj.GetNumeroDeCNPJ())
		user.SetFullName(cnpj.GetNomeDaPJ())
	} else {
		//validate the userCPF
		cpf, err1 := us.cpf.GetCpf(ctx, user.GetNationalID(), user.GetBornAt())
		if err1 != nil {
			err = err1
			return
		}
		user.SetNationalID(cpf.GetNumeroDeCpf())
		user.SetFullName(cpf.GetNomeDaPf())
	}

	//validate the user zipcode
	cep, err := us.globalService.GetCEP(ctx, user.GetZipCode())
	if err != nil {
		return
	}

	//validate the address number
	if user.GetNumber() == "" {
		err = utils.ErrInternalServer
		return
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
