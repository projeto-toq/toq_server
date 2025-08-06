package cache

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"google.golang.org/grpc"
)

func (c *cache) Get(ctx context.Context, fullMethod string, role usermodel.UserRole) (allowed bool, valid bool, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	service, method, err := c.DecodeFullmethod(fullMethod)
	privilege, valid := c.items[service].methodMap[method].rolemap[role]
	if !valid {
		return c.LoadNewPrivilege(ctx, service, method, role)
	}

	allowed = privilege.allowed
	valid = true

	return
}

func (c *cache) DecodeFullmethod(fullMethod string) (service usermodel.GRPCService, method uint8, err error) {
	paths := strings.Split(fullMethod, "/")
	if len(paths) < 3 {
		slog.Error("Error splinting fullMethod", "error", err)
		return 0, 0, errors.New("invalid full method")
	}

	switch paths[1] {
	case "grpc.UserService":
		method = c.GetMethodId(pb.UserService_ServiceDesc.Methods, paths[2])
		service = usermodel.ServiceUserService
	case "grpc.ListingService":
		method = c.GetMethodId(pb.ListingService_ServiceDesc.Methods, paths[2])
		service = usermodel.ServiceListingService
	}

	return
}

func (c *cache) GetMethodId(methods []grpc.MethodDesc, name string) uint8 {
	for i, method := range methods {
		if method.MethodName == name {
			return uint8(i)
		}
	}
	return 0
}

func (c *cache) LoadNewPrivilege(ctx context.Context, service usermodel.GRPCService, method uint8, role usermodel.UserRole) (allowed bool, valid bool, err error) {
	privilege, err1 := c.globalService.GetPrivilegeForCache(ctx, service, method, role)
	if err1 != nil {
		return false, false, err1
	}
	if privilege == nil {
		return false, true, nil
	}

	if _, ok := c.items[service]; !ok {
		methodStruct := MethodStruct{
			methodMap: make(map[uint8]RoleStruct),
		}
		c.items[service] = methodStruct
	}

	if _, ok := c.items[service].methodMap[method]; !ok {
		c.items[service].methodMap[method] = RoleStruct{
			rolemap: make(map[usermodel.UserRole]PrivilegeStruct),
		}
	}

	c.items[service].methodMap[method].rolemap[role] = PrivilegeStruct{
		privilege.Allowed(),
		time.Now().UTC(),
	}

	return privilege.Allowed(), true, nil
}
