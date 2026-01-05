package mediaprocessingservice

import (
	"context"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ReconcileStuckJobs marks running jobs that exceeded the timeout as failed and flags their assets.
func (s *mediaProcessingService) ReconcileStuckJobs(ctx context.Context, timeout time.Duration) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if timeout <= 0 {
		logger.Info("service.media.reconcile.disabled", "timeout", timeout.String())
		return nil
	}

	cutoff := s.nowUTC().Add(-timeout)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		return derrors.Infra("failed to start transaction", txErr)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("service.media.reconcile.rollback_error", "err", rbErr)
			}
		}
	}()

	stuckJobs, err := s.repo.ListStuckJobs(ctx, tx, mediaprocessingmodel.MediaProcessingJobStatusRunning, cutoff)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.reconcile.list_error", "err", err, "cutoff", cutoff)
		return derrors.Infra("failed to list stuck jobs", err)
	}

	if len(stuckJobs) == 0 {
		if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
			utils.SetSpanError(ctx, err)
			return derrors.Infra("failed to commit reconciliation transaction", err)
		}
		committed = true
		logger.Debug("service.media.reconcile.no_stuck_jobs", "cutoff", cutoff)
		return nil
	}

	now := s.nowUTC()
	for _, job := range stuckJobs {
		job.MarkCompleted(mediaprocessingmodel.MediaProcessingJobStatusFailed, mediaprocessingmodel.MediaProcessingJobPayload{}, now)
		job.AppendError(fmt.Sprintf("reconciler marked failed after %s without callback", timeout.String()))

		if updateErr := s.repo.UpdateProcessingJob(ctx, tx, job); updateErr != nil {
			utils.SetSpanError(ctx, updateErr)
			logger.Error("service.media.reconcile.update_job_error", "err", updateErr, "job_id", job.ID())
			continue
		}

		if err := s.repo.BulkUpdateAssetStatus(ctx, tx, job.ListingIdentityID(), mediaprocessingmodel.MediaAssetStatusProcessing, mediaprocessingmodel.MediaAssetStatusFailed); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("service.media.reconcile.bulk_fail_assets_error", "err", err, "listing_identity_id", job.ListingIdentityID())
		}
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.reconcile.commit_error", "err", err)
		return derrors.Infra("failed to commit reconciliation transaction", err)
	}
	committed = true

	logger.Warn("service.media.reconcile.completed", "stuck_jobs", len(stuckJobs), "cutoff", cutoff, "timeout", timeout)
	return nil
}
