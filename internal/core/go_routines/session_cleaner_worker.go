package goroutines

import (
	"context"
	"log/slog"
	"sync"
	"time"

	sessionservice "github.com/giulio-alfieri/toq_server/internal/core/service/session_service"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// SessionCleaner periodically deletes expired or fully expired sessions.
func SessionCleaner(svc sessionservice.Service, wg *sync.WaitGroup, ctx context.Context, interval time.Duration) {
	slog.Info("session cleaner routine started")
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			slog.Info("session cleaner routine stopped")
			return
		case <-ticker.C:
			// Do not generate traces for the cleaner routine
			noTraceCtx := utils.WithSkipTracing(ctx)
			if _, err := svc.CleanExpiredSessions(noTraceCtx, 500); err != nil {
				slog.Warn("session cleaner delete failed", "err", err)
			}
		}
	}
}
