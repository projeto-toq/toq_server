package goroutines

import (
	"context"
	"sync"
	"time"

	sessionservice "github.com/projeto-toq/toq_server/internal/core/service/session_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// SessionCleaner periodically deletes expired sessions with an additional retention cutoff.
func SessionCleaner(
	svc sessionservice.Service,
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

	if interval <= 0 {
		interval = time.Minute
	}
	if maxAge <= 0 {
		maxAge = 7 * 24 * time.Hour
	}
	if batchSize <= 0 {
		batchSize = 500
	}

	logger.Info("session cleaner routine started", "interval", interval, "max_age", maxAge, "batch_size", batchSize)

	runOnce := func(runCtx context.Context) {
		noTraceCtx := coreutils.WithSkipTracing(runCtx)
		cutoff := time.Now().Add(-maxAge)
		if _, err := svc.CleanExpiredSessionsBefore(noTraceCtx, cutoff, batchSize); err != nil {
			logger.Warn("session cleaner delete failed", "err", err)
		}
	}

	runOnce(ctx)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("session cleaner routine stopped")
			return
		case <-ticker.C:
			runOnce(ctx)
		}
	}
}
