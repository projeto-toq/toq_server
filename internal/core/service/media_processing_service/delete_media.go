package mediaprocessingservice

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
	"github.com/projeto-toq/toq_server/internal/core/domain/dto"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteMedia removes a media asset and all associated S3 objects in a synchronous, all-or-nothing flow.
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
		utils.SetSpanError(ctx, txErr)
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

	asset, err := s.repo.GetAsset(ctx, tx, uint64(input.ListingIdentityID), input.AssetType, input.Sequence)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.NotFound("asset not found")
		}
		utils.SetSpanError(ctx, err)
		return derrors.Infra("failed to get asset", err)
	}

	keys := collectDeletionKeys(asset)
	if len(keys) > 0 {
		if err := s.storage.DeleteKeys(ctx, keys); err != nil {
			utils.SetSpanError(ctx, err)
			utils.LoggerFromContext(ctx).Error("service.media.delete.s3_error", "err", err, "listing_identity_id", input.ListingIdentityID, "asset_type", input.AssetType, "sequence", input.Sequence)
			return err
		}
	}

	if err := s.repo.DeleteAsset(ctx, tx, uint64(input.ListingIdentityID), input.AssetType, input.Sequence); err != nil {
		utils.SetSpanError(ctx, err)
		return derrors.Infra("failed to delete asset from db", err)
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		return derrors.Infra("failed to commit transaction", err)
	}
	committed = true

	return nil
}

func collectDeletionKeys(asset mediaprocessingmodel.MediaAsset) []string {
	return dedupeDeletionKeys(asset.GetAllS3Keys())
}

func dedupeDeletionKeys(keys []string) []string {
	seen := make(map[string]struct{}, len(keys))
	ordered := make([]string, 0, len(keys))
	for _, key := range keys {
		trimmed := strings.TrimSpace(key)
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; exists {
			continue
		}
		seen[trimmed] = struct{}{}
		ordered = append(ordered, trimmed)
	}
	return ordered
}
