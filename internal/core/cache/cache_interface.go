package cache

import (
	"context"
	"time"

	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
)

// CacheInterface define a interface para operações de cache Redis
// Focado exclusivamente em cache de permissões de usuários
type CacheInterface interface {
	// Métodos de cache de permissões (principal funcionalidade)
	GetUserPermissions(ctx context.Context, userID int64) ([]byte, error)
	SetUserPermissions(ctx context.Context, userID int64, permissionsJSON []byte, ttl time.Duration) error
	DeleteUserPermissions(ctx context.Context, userID int64) error

	// Métodos de administração do cache
	Clean(ctx context.Context) // Limpeza geral do cache Redis
	Close() error              // Fechar conexão Redis

	// Injeção de dependências
	SetGlobalService(globalService globalservice.GlobalServiceInterface)
}
