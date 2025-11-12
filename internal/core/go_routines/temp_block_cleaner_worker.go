package goroutines

import (
	"context"
	"log/slog"
	"time"

	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	userservice "github.com/projeto-toq/toq_server/internal/core/service/user_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

type TempBlockCleanerWorker struct {
	userService   userservice.UserServiceInterface
	globalService globalservice.GlobalServiceInterface
	stopChan      chan struct{}
	doneChan      chan struct{}
	interval      time.Duration
	logger        *slog.Logger
}

func NewTempBlockCleanerWorker(
	userService userservice.UserServiceInterface,
	globalService globalservice.GlobalServiceInterface,
) *TempBlockCleanerWorker {
	return &TempBlockCleanerWorker{
		userService:   userService,
		globalService: globalService,
		stopChan:      make(chan struct{}),
		doneChan:      make(chan struct{}),
		interval:      5 * time.Minute, // Check every 5 minutes for expired blocks
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

	// Start transaction for batch query
	tx, err := w.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("Failed to start transaction for fetching expired blocks", "error", err)
		return
	}
	defer func() {
		if err != nil {
			if rbErr := w.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("Failed to rollback tx after fetch error", "error", rbErr)
			}
		}
	}()

	// Fetch users with expired blocks (returns []UserInterface from users table)
	expiredUsers, err := w.userService.GetUsersWithExpiredBlock(ctx, tx)
	if err != nil {
		logger.Error("Failed to get users with expired blocks", "error", err)
		return
	}

	// Commit the read transaction
	if cmErr := w.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		logger.Error("Failed to commit tx after fetching expired blocks", "error", cmErr)
		return
	}

	if len(expiredUsers) == 0 {
		logger.Debug("No expired temporary blocks found")
		return
	}

	logger.Info("Found users with expired blocks", "count", len(expiredUsers))

	for _, user := range expiredUsers {
		err := w.unblockUser(ctx, user.GetID())
		if err != nil {
			logger.Error("Failed to clear expired block", "userID", user.GetID(), "error", err)
			continue
		}
		logger.Info("User block cleared successfully", "userID", user.GetID())
	}
}

func (w *TempBlockCleanerWorker) unblockUser(ctx context.Context, userID int64) error {
	ctx = coreutils.ContextWithLogger(ctx)
	logger := coreutils.LoggerFromContext(ctx)

	// Start a new transaction for each user to avoid blocking other operations
	tx, err := w.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("Failed to start transaction for clearing block", "userID", userID, "error", err)
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := w.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("Failed to rollback tx when clearing block", "userID", userID, "error", rbErr)
			}
		}
	}()

	err = w.userService.ClearUserBlockedUntil(ctx, tx, userID)
	if err != nil {
		logger.Error("Failed to clear user blocked_until", "userID", userID, "error", err)
		return err
	}

	if cmErr := w.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		logger.Error("Failed to commit tx when clearing block", "userID", userID, "error", cmErr)
		return cmErr
	}

	return nil
}
