package goroutines

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/giulio-alfieri/toq_server/internal/core/cache"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
)

func CleanMemoryCache(memCache *cache.CacheInterface, wg *sync.WaitGroup, ctx context.Context) {
	slog.Info("memory cache cleaner routine started")
	ticker := time.NewTicker(globalmodel.ElapseTime)
	defer ticker.Stop()
	defer wg.Done()

	cache := *memCache

	for {
		select {
		case <-ctx.Done():
			slog.Info("memory cache cleaner routine stopped")
			return
		case <-ticker.C:
			cache.Clean(ctx)
			// slog.Info("memory cache cleaned")
		}
	}
}
