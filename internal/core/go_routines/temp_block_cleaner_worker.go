package goroutines

import (
	"context"
	"log/slog"
	"time"

	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	permissionservice "github.com/giulio-alfieri/toq_server/internal/core/service/permission_service"
)

type TempBlockCleanerWorker struct {
	permissionService permissionservice.PermissionServiceInterface
	globalService     globalservice.GlobalServiceInterface
	stopChan          chan struct{}
	doneChan          chan struct{}
	interval          time.Duration
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
	slog.Info("TempBlockCleanerWorker started")

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Process expired blocks immediately on startup
	w.processExpiredBlocks(ctx)

	go func() {
		defer close(w.doneChan)

		for {
			select {
			case <-ctx.Done():
				slog.Info("TempBlockCleanerWorker stopping due to context cancellation")
				return
			case <-w.stopChan:
				slog.Info("TempBlockCleanerWorker stopping due to stop signal")
				return
			case <-ticker.C:
				w.processExpiredBlocks(ctx)
			}
		}
	}()
}

func (w *TempBlockCleanerWorker) Stop() {
	slog.Info("TempBlockCleanerWorker stop requested")
	close(w.stopChan)
	<-w.doneChan
	slog.Info("TempBlockCleanerWorker stopped")
}

func (w *TempBlockCleanerWorker) processExpiredBlocks(ctx context.Context) {
	slog.Debug("Processing expired temporary blocks")

	expiredUsers, err := w.permissionService.GetExpiredTempBlockedUsers(ctx)
	if err != nil {
		slog.Error("Failed to get expired temp blocked users", "error", err)
		return
	}

	if len(expiredUsers) == 0 {
		slog.Debug("No expired temporary blocks found")
		return
	}

	slog.Info("Found expired temporary blocks", "count", len(expiredUsers))

	for _, userRole := range expiredUsers {
		err := w.unblockUser(ctx, userRole.GetUserID())
		if err != nil {
			slog.Error("Failed to unblock user", "userID", userRole.GetUserID(), "error", err)
			continue
		}
		slog.Info("User unblocked successfully", "userID", userRole.GetUserID())
	}
}

func (w *TempBlockCleanerWorker) unblockUser(ctx context.Context, userID int64) error {
	// Start a new transaction for each user to avoid blocking other operations
	tx, err := w.globalService.StartTransaction(ctx)
	if err != nil {
		slog.Error("Failed to start transaction for unblocking user", "userID", userID, "error", err)
		return err
	}
	defer func() {
		if err != nil {
			w.globalService.RollbackTransaction(ctx, tx)
		} else {
			w.globalService.CommitTransaction(ctx, tx)
		}
	}()

	err = w.permissionService.UnblockUser(ctx, tx, userID)
	if err != nil {
		slog.Error("Failed to unblock user in permission service", "userID", userID, "error", err)
		return err
	}

	return nil
}
