package cache

import (
	"context"
	"time"
)

func (c *cache) Clean(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for service, methodStruct := range c.items {
		c.CheckExpiredMethod(ctx, methodStruct)
		if len(methodStruct.methodMap) == 0 {
			delete(c.items, service)
		}
	}
}

func (c *cache) CheckExpiredMethod(ctx context.Context, methodStruct MethodStruct) {

	for methodID, roleStruct := range methodStruct.methodMap {
		c.CheckExpiredRole(ctx, roleStruct)
		if len(roleStruct.rolemap) == 0 {
			delete(methodStruct.methodMap, methodID)
		}
	}
}

func (c *cache) CheckExpiredRole(ctx context.Context, roleStruct RoleStruct) {
	for role, privilege := range roleStruct.rolemap {
		if time.Since(privilege.lastAccess) > time.Duration(1*time.Minute) {
			delete(roleStruct.rolemap, role)
		}
	}

}
