package goroutines

import (
	"context"
	"time"

	validationservice "github.com/giulio-alfieri/toq_server/internal/core/service/validation_service"
	coreutils "github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// ValidationCleaner periodically deletes expired rows from temp_user_validations.
func ValidationCleaner(svc validationservice.Service, interval time.Duration, ctx context.Context) {
	ctx = coreutils.ContextWithLogger(ctx)
	logger := coreutils.LoggerFromContext(ctx)

	logger.Info("validation cleaner worker started", "interval", interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Kick off an immediate cleanup on start
	if _, err := svc.CleanExpiredValidations(coreutils.WithSkipTracing(ctx), 500); err != nil {
		logger.Warn("validation cleaner immediate run failed", "err", err)
	}

	for {
		select {
		case <-ctx.Done():
			logger.Info("validation cleaner worker stopped")
			return
		case <-ticker.C:
			noTraceCtx := coreutils.WithSkipTracing(ctx)
			if _, err := svc.CleanExpiredValidations(noTraceCtx, 500); err != nil {
				logger.Warn("validation cleaner delete failed", "err", err)
			}
		}
	}
}
