package goroutines

import (
	"context"
	"time"

	mediaprocessingservice "github.com/projeto-toq/toq_server/internal/core/service/media_processing_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// MediaProcessingReconciler periodically marks timed-out processing jobs as failed.
type MediaProcessingReconciler struct {
	service  mediaprocessingservice.MediaProcessingServiceInterface
	interval time.Duration
	timeout  time.Duration
}

// NewMediaProcessingReconciler builds a reconciler worker with sane defaults.
func NewMediaProcessingReconciler(service mediaprocessingservice.MediaProcessingServiceInterface, interval, timeout time.Duration) *MediaProcessingReconciler {
	if interval <= 0 {
		interval = 5 * time.Minute
	}
	return &MediaProcessingReconciler{service: service, interval: interval, timeout: timeout}
}

// Start runs the reconciliation loop until the context is cancelled.
func (w *MediaProcessingReconciler) Start(ctx context.Context) {
	logger := coreutils.LoggerFromContext(ctx)
	if w.service == nil {
		logger.Warn("goroutine.media_reconciler.service_nil")
		return
	}
	if w.timeout <= 0 {
		logger.Info("goroutine.media_reconciler.disabled", "timeout", w.timeout.String())
		return
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.runOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			logger.Info("goroutine.media_reconciler.stopped")
			return
		case <-ticker.C:
			w.runOnce(ctx)
		}
	}
}

func (w *MediaProcessingReconciler) runOnce(ctx context.Context) {
	logger := coreutils.LoggerFromContext(ctx)
	if err := w.service.ReconcileStuckJobs(ctx, w.timeout); err != nil {
		logger.Error("goroutine.media_reconciler.error", "err", err, "timeout", w.timeout.String())
	} else {
		logger.Debug("goroutine.media_reconciler.ok", "timeout", w.timeout.String())
	}
}
