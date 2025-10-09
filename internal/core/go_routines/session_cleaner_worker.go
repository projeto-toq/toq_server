package goroutines

import (
	"context"
	"sync"
	"time"

	sessionservice "github.com/projeto-toq/toq_server/internal/core/service/session_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// SessionCleaner periodically deletes expired or fully expired sessions.
func SessionCleaner(svc sessionservice.Service, wg *sync.WaitGroup, ctx context.Context, interval time.Duration) {
	ctx = coreutils.ContextWithLogger(ctx)
	logger := coreutils.LoggerFromContext(ctx)

	logger.Info("session cleaner routine started")
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			logger.Info("session cleaner routine stopped")
			return
		case <-ticker.C:
			// Do not generate traces for the cleaner routine
			noTraceCtx := coreutils.WithSkipTracing(ctx)
			if _, err := svc.CleanExpiredSessions(noTraceCtx, 500); err != nil {
				logger.Warn("session cleaner delete failed", "err", err)
			}
		}
	}
}
