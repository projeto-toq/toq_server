package cache

import (
	"context"
	"sync"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
)

type CacheInterface interface {
	Get(ctx context.Context, fullMethod string, role usermodel.UserRole) (allowed bool, valid bool, err error)
	Clean(ctx context.Context)
	Close() error
	SetGlobalService(globalService globalservice.GlobalServiceInterface) // Para injeção posterior
}

var instance *cache
var once sync.Once

func NewCache(globalService globalservice.GlobalServiceInterface) CacheInterface {
	once.Do(func() {
		instance = &cache{
			items:         make(map[usermodel.GRPCService](MethodStruct)),
			globalService: globalService,
		}
	})
	return instance
}
