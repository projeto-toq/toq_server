package goroutines

import (
	"context"
	"sync"
	"time"

	userservices "github.com/projeto-toq/toq_server/internal/core/service/user_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeviceTokenCleaner periodically purges stale device tokens according to retention policies.
func DeviceTokenCleaner(
	svc userservices.UserServiceInterface,
	wg *sync.WaitGroup,
	ctx context.Context,
	interval time.Duration,
	maxAge time.Duration,
	batchSize int,
) {
	ctx = coreutils.ContextWithLogger(ctx)
	logger := coreutils.LoggerFromContext(ctx)

	if wg != nil {
		defer wg.Done()
	}

	if svc == nil {
		logger.Warn("device_token cleaner skipped: service unavailable")
		return
	}

	if interval <= 0 {
		interval = time.Hour
	}
	if maxAge <= 0 {
		maxAge = 90 * 24 * time.Hour
	}
	if batchSize <= 0 {
		batchSize = 500
	}

	logger.Info("device_token cleaner started", "interval", interval, "max_age", maxAge, "batch_size", batchSize)

	runOnce := func(runCtx context.Context) {
		noTraceCtx := coreutils.WithSkipTracing(runCtx)
		deleted, err := svc.PurgeStaleDeviceTokens(noTraceCtx, maxAge, batchSize)
		if err != nil {
			logger.Warn("device_token.cleaner.run_failed", "err", err)
			return
		}
		if deleted > 0 {
			logger.Info("device_token.cleaner.deleted", "count", deleted)
		}
	}

	runOnce(ctx)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("device_token cleaner stopped")
			return
		case <-ticker.C:
			runOnce(ctx)
		}
	}
}
