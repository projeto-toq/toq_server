package userservices

import (
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"google.golang.org/grpc"
)

func (us *userService) CreateRealtorPrivileges() (privileges []usermodel.PrivilegeInterface) {
	us.AddUserRealtorUserPrivileges(pb.UserService_ServiceDesc.Methods, &privileges)
	us.AddListingRealtorUserPrivileges(pb.UserService_ServiceDesc.Methods, &privileges)
	return
}

func (us *userService) AddUserRealtorUserPrivileges(methods []grpc.MethodDesc, privileges *[]usermodel.PrivilegeInterface) {

	for methodID, method := range methods {
		privilege := usermodel.NewPrivilege()
		privilege.SetService(usermodel.ServiceUserService)
		privilege.SetMethod(uint8(methodID))
		privilege.SetAllowed(usermodel.RealtorUserPrivileges[method.MethodName])
		*privileges = append(*privileges, privilege)
	}
}

func (us *userService) AddListingRealtorUserPrivileges(methods []grpc.MethodDesc, privileges *[]usermodel.PrivilegeInterface) {

	for methodID, method := range methods {
		privilege := usermodel.NewPrivilege()
		privilege.SetService(usermodel.ServiceListingService)
		privilege.SetMethod(uint8(methodID))
		privilege.SetAllowed(usermodel.RealtorListingPrivileges[method.MethodName])
		*privileges = append(*privileges, privilege)
	}
}
