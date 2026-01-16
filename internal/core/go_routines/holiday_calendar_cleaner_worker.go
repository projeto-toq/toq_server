package goroutines

import (
	"context"
	"sync"
	"time"

	holidayservices "github.com/projeto-toq/toq_server/internal/core/service/holiday_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// HolidayCalendarCleaner periodically deletes old non-recurrent holiday dates based on retention policy.
func HolidayCalendarCleaner(
	svc holidayservices.HolidayServiceInterface,
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
		logger.Warn("holiday cleaner skipped: service unavailable")
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

	logger.Info("holiday cleaner started", "interval", interval, "max_age", maxAge, "batch_size", batchSize)

	runOnce := func(runCtx context.Context) {
		cutoff := time.Now().Add(-maxAge)
		noTraceCtx := coreutils.WithSkipTracing(runCtx)
		deleted, err := svc.CleanOldCalendarDates(noTraceCtx, cutoff, batchSize)
		if err != nil {
			logger.Warn("holiday.cleaner.run_failed", "err", err)
			return
		}
		if deleted > 0 {
			logger.Info("holiday.cleaner.deleted", "count", deleted)
		}
	}

	runOnce(ctx)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("holiday cleaner stopped")
			return
		case <-ticker.C:
			runOnce(ctx)
		}
	}
}
