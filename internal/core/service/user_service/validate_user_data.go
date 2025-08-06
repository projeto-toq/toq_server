package userservices

import (
	"context"
	"database/sql"
	"time"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (us *userService) ValidateUserData(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface, role usermodel.UserRole) (err error) {

	now := time.Now().UTC()

	//verify if user already exists
	exist, err := us.repo.VerifyUserDuplicity(ctx, tx, user)
	if err != nil {
		return
	}

	if exist {
		err = status.Error(codes.AlreadyExists, "User already exists")
		return
	}

	//verify the password
	err = validatePassword(user.GetPassword())
	if err != nil {
		return
	}

	if role == usermodel.RoleAgency {

		cnpj, err1 := us.cnpj.GetCNPJ(ctx, user.GetNationalID()) //TODO: mover para o global service, assim como cep e cpf
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
		err = status.Error(codes.InvalidArgument, "Address number is required")
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
