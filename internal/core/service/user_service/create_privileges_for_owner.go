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
	for methodID := range methods {
		privilege := usermodel.NewPrivilege()
		privilege.SetService(usermodel.ServiceUserService)
		privilege.SetMethod(uint8(methodID))
		// TODO: Replace with HTTP privileges after migration
		// privilege.SetAllowed(usermodel.OwnerUserPrivileges[method.MethodName])
		privilege.SetAllowed(false) // Temporarily disabled during HTTP migration
		*privileges = append(*privileges, privilege)
	}
}

func (us *userService) AddOwnerListingPrivileges(methods []grpc.MethodDesc, privileges *[]usermodel.PrivilegeInterface) {
	for methodID := range methods {
		privilege := usermodel.NewPrivilege()
		privilege.SetService(usermodel.ServiceListingService)
		privilege.SetMethod(uint8(methodID))
		// TODO: Replace with HTTP privileges after migration
		// privilege.SetAllowed(usermodel.OwnerListingPrivileges[method.MethodName])
		privilege.SetAllowed(false) // Temporarily disabled during HTTP migration
		*privileges = append(*privileges, privilege)
	}
}
