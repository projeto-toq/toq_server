package mediaprocessingservice

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	"github.com/projeto-toq/toq_server/internal/core/domain/dto"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateMedia modifies metadata or properties of an existing asset.
func (s *mediaProcessingService) UpdateMedia(ctx context.Context, input dto.UpdateMediaInput) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	if input.ListingIdentityID == 0 {
		return derrors.Validation("listingIdentityId must be greater than zero", map[string]any{"listingIdentityId": "required"})
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		return derrors.Infra("failed to start transaction", txErr)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.LoggerFromContext(ctx).Error("service.media.update.rollback_error", "err", rbErr)
			}
		}
	}()

	// Find the asset
	asset, err := s.repo.GetAsset(ctx, tx, uint64(input.ListingIdentityID), input.AssetType, input.Sequence)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.NotFound("asset not found")
		}
		return derrors.Infra("failed to get asset", err)
	}

	// Update fields
	if input.Title != "" {
		asset.SetTitle(input.Title)
	}
	if len(input.Metadata) > 0 {
		currentMeta := make(map[string]string)
		if asset.Metadata() != "" {
			_ = json.Unmarshal([]byte(asset.Metadata()), &currentMeta)
		}
		for k, v := range input.Metadata {
			currentMeta[k] = v
		}
		metaBytes, _ := json.Marshal(currentMeta)
		asset.SetMetadata(string(metaBytes))
	}

	if err := s.repo.UpsertAsset(ctx, tx, asset); err != nil {
		return derrors.Infra("failed to update asset", err)
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		return derrors.Infra("failed to commit transaction", err)
	}
	committed = true

	return nil
}

// DeleteMedia removes an asset from the database and storage.
func (s *mediaProcessingService) DeleteMedia(ctx context.Context, input dto.DeleteMediaInput) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	if input.ListingIdentityID == 0 {
		return derrors.Validation("listingIdentityId must be greater than zero", map[string]any{"listingIdentityId": "required"})
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		return derrors.Infra("failed to start transaction", txErr)
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.LoggerFromContext(ctx).Error("service.media.delete.rollback_error", "err", rbErr)
			}
		}
	}()

	// Check if asset exists
	asset, err := s.repo.GetAsset(ctx, tx, uint64(input.ListingIdentityID), input.AssetType, input.Sequence)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.NotFound("asset not found")
		}
		return derrors.Infra("failed to get asset", err)
	}

	// Delete from DB first (transactional)
	if err := s.repo.DeleteAsset(ctx, tx, uint64(input.ListingIdentityID), input.AssetType, input.Sequence); err != nil {
		return derrors.Infra("failed to delete asset from db", err)
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		return derrors.Infra("failed to commit transaction", err)
	}
	committed = true

	// Cleanup Storage (Best effort, after commit)
	// Deletes all associated files (Raw, Processed, Thumbnails, etc.)
	go func() {
		bgCtx := context.Background() // Use background context for async cleanup
		keysToDelete := asset.GetAllS3Keys()

		for _, key := range keysToDelete {
			if err := s.storage.DeleteObject(bgCtx, key); err != nil {
				// Log warning but continue trying to delete other keys
				utils.LoggerFromContext(ctx).Warn("service.media.delete.s3_failed", "key", key, "err", err)
			}
		}
	}()

	return nil
}
