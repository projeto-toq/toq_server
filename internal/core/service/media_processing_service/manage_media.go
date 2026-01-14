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
