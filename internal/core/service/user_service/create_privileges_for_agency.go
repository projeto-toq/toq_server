package userservices

import (
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"google.golang.org/grpc"
)

func (us *userService) CreateAgencyPrivileges() (privileges []usermodel.PrivilegeInterface) {
	us.AddAgencyUserPrivileges(pb.UserService_ServiceDesc.Methods, &privileges)
	return
}

func (us *userService) AddAgencyUserPrivileges(methods []grpc.MethodDesc, privileges *[]usermodel.PrivilegeInterface) {
	for methodID := range methods {
		privilege := usermodel.NewPrivilege()
		privilege.SetService(usermodel.ServiceUserService)
		privilege.SetMethod(uint8(methodID))
		// TODO: Replace with HTTP privileges after migration
		// privilege.SetAllowed(usermodel.AgencyUserPrivileges[method.MethodName])
		privilege.SetAllowed(false) // Temporarily disabled during HTTP migration
		*privileges = append(*privileges, privilege)
	}
}
