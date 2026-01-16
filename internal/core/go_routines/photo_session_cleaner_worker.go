package goroutines

import (
	"context"
	"sync"
	"time"

	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PhotoSessionCleaner periodically removes old bookings and orphan agenda entries per retention policy.
func PhotoSessionCleaner(
	svc photosessionservices.PhotoSessionServiceInterface,
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
		logger.Warn("photo_session cleaner skipped: service unavailable")
		return
	}

	if interval <= 0 {
		interval = 12 * time.Hour
	}
	if maxAge <= 0 {
		maxAge = 365 * 24 * time.Hour
	}
	if batchSize <= 0 {
		batchSize = 500
	}

	logger.Info("photo_session cleaner started", "interval", interval, "max_age", maxAge, "batch_size", batchSize)

	runOnce := func(runCtx context.Context) {
		cutoff := time.Now().Add(-maxAge)
		noTraceCtx := coreutils.WithSkipTracing(runCtx)

		if deleted, err := svc.CleanOldBookings(noTraceCtx, cutoff, batchSize); err != nil {
			logger.Warn("photo_session.cleaner.bookings_failed", "err", err)
		} else if deleted > 0 {
			logger.Info("photo_session.cleaner.bookings_deleted", "count", deleted)
		}

		if deleted, err := svc.CleanOldAgendaEntries(noTraceCtx, cutoff, batchSize); err != nil {
			logger.Warn("photo_session.cleaner.agenda_failed", "err", err)
		} else if deleted > 0 {
			logger.Info("photo_session.cleaner.agenda_deleted", "count", deleted)
		}
	}

	runOnce(ctx)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("photo_session cleaner stopped")
			return
		case <-ticker.C:
			runOnce(ctx)
		}
	}
}
