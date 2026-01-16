package goroutines

import (
	"context"
	"sync"
	"time"

	mediaprocessingservice "github.com/projeto-toq/toq_server/internal/core/service/media_processing_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// MediaProcessingCleaner periodically removes terminal media processing jobs older than retention window.
func MediaProcessingCleaner(
	svc mediaprocessingservice.MediaProcessingServiceInterface,
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
		logger.Warn("media_processing cleaner skipped: service unavailable")
		return
	}

	if interval <= 0 {
		interval = 2 * time.Hour
	}
	if maxAge <= 0 {
		maxAge = 90 * 24 * time.Hour
	}
	if batchSize <= 0 {
		batchSize = 500
	}

	logger.Info("media_processing cleaner started", "interval", interval, "max_age", maxAge, "batch_size", batchSize)

	runOnce := func(runCtx context.Context) {
		cutoff := time.Now().Add(-maxAge)
		noTraceCtx := coreutils.WithSkipTracing(runCtx)
		deleted, err := svc.CleanOldJobs(noTraceCtx, cutoff, batchSize)
		if err != nil {
			logger.Warn("media_processing.cleaner.run_failed", "err", err)
			return
		}
		if deleted > 0 {
			logger.Info("media_processing.cleaner.deleted", "count", deleted)
		}
	}

	runOnce(ctx)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("media_processing cleaner stopped")
			return
		case <-ticker.C:
			runOnce(ctx)
		}
	}
}
