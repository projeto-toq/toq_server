package userHandlerconverters

import (
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

func UserDomainToProfileResponse(user usermodel.UserInterface) (response *pb.GetProfileResponse, err error) {

	role := pb.UserRole{}
	role.Id = user.GetActiveRole().GetID()
	role.Role = usermodel.UserRole.String(user.GetActiveRole().GetRole())
	role.Active = user.GetActiveRole().IsActive()
	role.Status = usermodel.UserRoleStatus.String(user.GetActiveRole().GetStatus())
	role.StatusReason = user.GetActiveRole().GetStatusReason()

	userResp := &pb.User{
		Id:            user.GetID(),
		ActiveRole:    &role,
		FullName:      user.GetFullName(),
		NickName:      user.GetNickName(),
		NationalID:    user.GetNationalID(),
		CreciNumber:   user.GetCreciNumber(),
		CreciState:    user.GetCreciState(),
		CreciValidity: user.GetCreciValidity().Format("2006-01-02"),
		BornAT:        user.GetBornAt().Format("2006-01-02"),
		PhoneNumber:   user.GetPhoneNumber(),
		Email:         user.GetEmail(),
		ZipCode:       user.GetZipCode(),
		Street:        user.GetStreet(),
		Number:        user.GetNumber(),
		Complement:    user.GetComplement(),
		Neighborhood:  user.GetNeighborhood(),
		City:          user.GetCity(),
		State:         user.GetState(),
	}

	response = &pb.GetProfileResponse{User: userResp}

	return
}
