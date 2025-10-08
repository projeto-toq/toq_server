package goroutines

import (
	"context"
	"log/slog"
	"time"

	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	permissionservice "github.com/giulio-alfieri/toq_server/internal/core/service/permission_service"
	coreutils "github.com/giulio-alfieri/toq_server/internal/core/utils"
)

type TempBlockCleanerWorker struct {
	permissionService permissionservice.PermissionServiceInterface
	globalService     globalservice.GlobalServiceInterface
	stopChan          chan struct{}
	doneChan          chan struct{}
	interval          time.Duration
	logger            *slog.Logger
}

func NewTempBlockCleanerWorker(
	permissionService permissionservice.PermissionServiceInterface,
	globalService globalservice.GlobalServiceInterface,
) *TempBlockCleanerWorker {
	return &TempBlockCleanerWorker{
		permissionService: permissionService,
		globalService:     globalService,
		stopChan:          make(chan struct{}),
		doneChan:          make(chan struct{}),
		interval:          5 * time.Minute, // Check every 5 minutes for expired blocks
	}
}

func (w *TempBlockCleanerWorker) Start(ctx context.Context) {
	ctx = coreutils.ContextWithLogger(ctx)
	w.logger = coreutils.LoggerFromContext(ctx)
	w.logger.Info("TempBlockCleanerWorker started")

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Process expired blocks immediately on startup
	w.processExpiredBlocks(ctx)

	go func(loopCtx context.Context) {
		loopCtx = coreutils.ContextWithLogger(loopCtx)
		defer close(w.doneChan)

		for {
			select {
			case <-loopCtx.Done():
				w.logger.Info("TempBlockCleanerWorker stopping due to context cancellation")
				return
			case <-w.stopChan:
				w.logger.Info("TempBlockCleanerWorker stopping due to stop signal")
				return
			case <-ticker.C:
				w.processExpiredBlocks(loopCtx)
			}
		}
	}(ctx)
}

func (w *TempBlockCleanerWorker) Stop() {
	logger := w.logger
	if logger == nil {
		logger = slog.Default()
	}
	logger.Info("TempBlockCleanerWorker stop requested")
	close(w.stopChan)
	<-w.doneChan
	logger.Info("TempBlockCleanerWorker stopped")
}

func (w *TempBlockCleanerWorker) processExpiredBlocks(ctx context.Context) {
	ctx = coreutils.ContextWithLogger(ctx)
	logger := coreutils.LoggerFromContext(ctx)

	logger.Debug("Processing expired temporary blocks")

	expiredUsers, err := w.permissionService.GetExpiredTempBlockedUsers(ctx)
	if err != nil {
		logger.Error("Failed to get expired temp blocked users", "error", err)
		return
	}

	if len(expiredUsers) == 0 {
		logger.Debug("No expired temporary blocks found")
		return
	}

	logger.Info("Found expired temporary blocks", "count", len(expiredUsers))

	for _, userRole := range expiredUsers {
		err := w.unblockUser(ctx, userRole.GetUserID())
		if err != nil {
			logger.Error("Failed to unblock user", "userID", userRole.GetUserID(), "error", err)
			continue
		}
		logger.Info("User unblocked successfully", "userID", userRole.GetUserID())
	}
}

func (w *TempBlockCleanerWorker) unblockUser(ctx context.Context, userID int64) error {
	ctx = coreutils.ContextWithLogger(ctx)
	logger := coreutils.LoggerFromContext(ctx)

	// Start a new transaction for each user to avoid blocking other operations
	tx, err := w.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("Failed to start transaction for unblocking user", "userID", userID, "error", err)
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := w.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("Failed to rollback tx when unblocking user", "userID", userID, "error", rbErr)
			}
		}
	}()

	err = w.permissionService.UnblockUser(ctx, tx, userID)
	if err != nil {
		logger.Error("Failed to unblock user in permission service", "userID", userID, "error", err)
		return err
	}

	if cmErr := w.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		logger.Error("Failed to commit tx when unblocking user", "userID", userID, "error", cmErr)
		return cmErr
	}

	return nil
}
