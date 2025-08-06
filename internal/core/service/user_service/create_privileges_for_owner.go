package userservices

import (
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"google.golang.org/grpc"
)

func (us *userService) CreateOwnerPrivileges() (privileges []usermodel.PrivilegeInterface) {
	us.AddOwnerUserPrivileges(pb.UserService_ServiceDesc.Methods, &privileges)
	us.AddOwnerListingPrivileges(pb.ListingService_ServiceDesc.Methods, &privileges)
	return

}

func (us *userService) AddOwnerUserPrivileges(methods []grpc.MethodDesc, privileges *[]usermodel.PrivilegeInterface) {

	for methodID, method := range methods {
		privilege := usermodel.NewPrivilege()
		privilege.SetService(usermodel.ServiceUserService)
		privilege.SetMethod(uint8(methodID))
		privilege.SetAllowed(usermodel.OwnerUserPrivileges[method.MethodName])
		*privileges = append(*privileges, privilege)
	}
}

func (us *userService) AddOwnerListingPrivileges(methods []grpc.MethodDesc, privileges *[]usermodel.PrivilegeInterface) {
	for methodID, method := range methods {
		privilege := usermodel.NewPrivilege()
		privilege.SetService(usermodel.ServiceListingService)
		privilege.SetMethod(uint8(methodID))
		privilege.SetAllowed(usermodel.OwnerListingPrivileges[method.MethodName])
		*privileges = append(*privileges, privilege)
	}

}
