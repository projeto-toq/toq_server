package userservices

import (
	// TODO: Replace with HTTP-based privilege system
	// // "github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	// // 
)

func (us *userService) CreateAgencyPrivileges() (privileges []usermodel.PrivilegeInterface) {
	// TODO: Implement HTTP-based privilege creation
	// Temporarily creating basic privileges
	us.addBasicAgencyPrivileges(&privileges)
	return
}

// Temporary implementation during HTTP migration
func (us *userService) addBasicAgencyPrivileges(privileges *[]usermodel.PrivilegeInterface) {
	// Create basic privileges for agency users
	// This is a temporary implementation during HTTP migration
	for i := 0; i < 10; i++ { // Create 10 basic privilege entries
		privilege := usermodel.NewPrivilege()
		privilege.SetService(usermodel.ServiceUserService)
		privilege.SetMethod(uint8(i))
		privilege.SetAllowed(false) // Conservative default during migration
		*privileges = append(*privileges, privilege)
	}
}

// Legacy method - will be replaced with HTTP-based implementation
func (us *userService) AddAgencyUserPrivileges(methods interface{}, privileges *[]usermodel.PrivilegeInterface) {
	// Temporarily disabled during HTTP migration
	// TODO: Implement HTTP-based privilege assignment
}
