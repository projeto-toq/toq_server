package cache

import (
	"sync"
	"time"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
)

type cache struct {
	globalService globalservice.GlobalServiceInterface
	items         map[usermodel.GRPCService](MethodStruct)
	mu            sync.Mutex
}

type MethodStruct struct {
	methodMap map[uint8](RoleStruct)
}

type RoleStruct struct {
	rolemap map[usermodel.UserRole]PrivilegeStruct
}

type PrivilegeStruct struct {
	allowed    bool
	lastAccess time.Time
}
