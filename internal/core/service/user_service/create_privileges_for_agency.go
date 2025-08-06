package userservices

import (
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

func (us *userService) CreateAgencyPrivileges() (privileges []usermodel.PrivilegeInterface) {

	for methodID, method := range pb.UserService_ServiceDesc.Methods {
		privilege := usermodel.NewPrivilege()
		privilege.SetService(usermodel.ServiceUserService)
		privilege.SetMethod(uint8(methodID))
		privilege.SetAllowed(usermodel.AgencyUserPrivileges[method.MethodName])
		privileges = append(privileges, privilege)
	}
	return
}
