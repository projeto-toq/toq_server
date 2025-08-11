package cache

import (
	"log/slog"

	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
)

// SetGlobalService injeta o GlobalService após a criação do cache
// Usado para resolver dependências circulares entre Cache e GlobalService
func (c *cache) SetGlobalService(globalService globalservice.GlobalServiceInterface) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.globalService = globalService
	slog.Debug("GlobalService injected into cache")
}
