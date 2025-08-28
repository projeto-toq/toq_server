package cache

import (
	"context"
	"log/slog"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

// Temporary implementation of cache methods for HTTP migration
// TODO: Implement proper HTTP-based permission cache system

func (c *cache) Get(ctx context.Context, fullMethod string, role usermodel.UserRole) (allowed bool, valid bool, err error) {
	// TODO: Implement HTTP-based permission cache system
	// For now, return permissive defaults to maintain functionality
	slog.Debug("Cache system temporarily disabled - allowing all requests",
		"method", fullMethod, "role", role)
	return true, true, nil
}

// DecodeFullmethod temporarily disabled for HTTP migration
func (c *cache) DecodeFullmethod(fullMethod string) (service usermodel.GRPCService, method uint8, err error) {
	// Temporary implementation
	return usermodel.ServiceUserService, 0, nil
}

// GetMethodId temporarily disabled for HTTP migration
func (c *cache) GetMethodId(methods interface{}, name string) uint8 {
	// Temporary implementation
	return 0
}

// LoadNewPrivilege temporarily disabled for HTTP migration
func (c *cache) LoadNewPrivilege(ctx context.Context, service usermodel.GRPCService, method uint8, role usermodel.UserRole) (allowed bool, valid bool, err error) {
	// Temporary implementation - always allow for now
	slog.Debug("LoadNewPrivilege temporarily disabled - allowing request",
		"service", service, "method", method, "role", role)
	return true, true, nil
}
