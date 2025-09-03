package goroutines

import (
	"context"
	"log/slog"
	"time"

	validationservice "github.com/giulio-alfieri/toq_server/internal/core/service/validation_service"
)

// ValidationCleaner periodically deletes expired rows from temp_user_validations.
func ValidationCleaner(svc validationservice.Service, interval time.Duration, ctx context.Context) {
	slog.Info("validation cleaner worker started", "interval", interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Kick off an immediate cleanup on start
	if _, err := svc.CleanExpiredValidations(ctx, 500); err != nil {
		slog.Warn("validation cleaner immediate run failed", "err", err)
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("validation cleaner worker stopped")
			return
		case <-ticker.C:
			if _, err := svc.CleanExpiredValidations(ctx, 500); err != nil {
				slog.Warn("validation cleaner delete failed", "err", err)
			}
		}
	}
}
