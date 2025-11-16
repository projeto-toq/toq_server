package mediaprocessingservice

import (
	"context"
	"database/sql"
	"errors"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetBatchStatus retrieves detailed status for a specific batch under a listing identity.
func (s *mediaProcessingService) GetBatchStatus(ctx context.Context, input GetBatchStatusInput) (GetBatchStatusOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return GetBatchStatusOutput{}, derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.ListingID == 0 {
		return GetBatchStatusOutput{}, derrors.Validation("listingId must be greater than zero", map[string]any{"listingId": "required"})
	}
	if input.BatchID == 0 {
		return GetBatchStatusOutput{}, derrors.Validation("batchId must be greater than zero", map[string]any{"batchId": "required"})
	}

	tx, txErr := s.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("service.media.get_batch_status.tx_start_error", "err", txErr, "listing_id", input.ListingID)
		return GetBatchStatusOutput{}, derrors.Infra("failed to start transaction", txErr)
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("service.media.get_batch_status.tx_rollback_error", "err", rbErr, "listing_id", input.ListingID)
		}
	}()

	batch, err := s.repo.GetBatchByID(ctx, tx, input.BatchID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return GetBatchStatusOutput{}, derrors.NotFound("batch not found")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.get_batch_status.get_batch_error", "err", err, "batch_id", input.BatchID)
		return GetBatchStatusOutput{}, derrors.Infra("failed to load batch", err)
	}

	if batch.ListingID() != input.ListingID {
		return GetBatchStatusOutput{}, derrors.Conflict("batch does not belong to listing")
	}

	assets, err := s.repo.ListAssetsByBatch(ctx, tx, input.BatchID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("service.media.get_batch_status.list_assets_error", "err", err, "batch_id", input.BatchID)
		return GetBatchStatusOutput{}, derrors.Infra("failed to list assets", err)
	}

	assetStatuses := make([]BatchAssetStatus, 0, len(assets))
	for _, asset := range assets {
		clientID := asset.Metadata()["client_id"]
		title := asset.Metadata()["title"]

		assetStatuses = append(assetStatuses, BatchAssetStatus{
			ClientID:     clientID,
			Title:        title,
			AssetType:    asset.AssetType(),
			Sequence:     asset.Sequence(),
			RawObjectKey: asset.RawObjectKey(),
			ProcessedKey: asset.ProcessedKey(),
			ThumbnailKey: asset.ThumbnailKey(),
			Metadata:     cloneStringMap(asset.Metadata()),
		})
	}

	statusMetadata := batch.StatusMetadata()
	statusMessage := statusMetadata.Message
	if statusMetadata.Reason != "" {
		statusMessage = statusMetadata.Reason
	}

	logger.Info("service.media.get_batch_status.success",
		"listing_id", input.ListingID,
		"batch_id", input.BatchID,
		"status", batch.Status(),
		"assets", len(assets),
	)

	return GetBatchStatusOutput{
		ListingID:     input.ListingID,
		BatchID:       input.BatchID,
		Status:        batch.Status(),
		StatusMessage: statusMessage,
		Assets:        assetStatuses,
	}, nil
}
