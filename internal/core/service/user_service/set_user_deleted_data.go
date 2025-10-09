package userservices

import (
	"time"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

func (us *userService) setDeletedData(user usermodel.UserInterface) {

	// Future enhancement: Check and delete associated listings for property owners
	// Status do user_role é marcado como StatusDeleted em markUserRolesAsDeleted
	user.SetFullName("Apagado por solicitação do usuário")
	user.SetNickName("Apagado")
	user.SetNationalID("00000000000")
	user.SetCreciNumber("000000")
	user.SetCreciState(" ")
	user.SetCreciValidity(time.Now().UTC())
	user.SetBornAt(time.Now().UTC())
	user.SetPhoneNumber("+000000000000")
	user.SetEmail("")
	user.SetZipCode("")
	user.SetStreet("")
	user.SetNumber("")
	user.SetComplement("")
	user.SetNeighborhood("")
	user.SetCity("")
	user.SetState("")
	user.SetPassword("")
	user.SetDeleted(true)
}
