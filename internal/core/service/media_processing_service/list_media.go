package mediaprocessingservice

import (
	"context"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	"github.com/projeto-toq/toq_server/internal/core/domain/dto"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	mediaprocessingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/media_processing_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *mediaProcessingService) ListMedia(ctx context.Context, input dto.ListMediaInput) (dto.ListMediaOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return dto.ListMediaOutput{}, derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	if input.ListingIdentityID == 0 {
		return dto.ListMediaOutput{}, derrors.Validation("listingIdentityId required", nil)
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		return dto.ListMediaOutput{}, derrors.Infra("failed to start transaction", txErr)
	}
	defer func() { _ = s.globalService.RollbackTransaction(ctx, tx) }()

	filter := mediaprocessingrepository.AssetFilter{
		Sequence: input.Sequence,
	}
	if input.AssetType != "" {
		filter.AssetTypes = []mediaprocessingmodel.MediaAssetType{mediaprocessingmodel.MediaAssetType(input.AssetType)}
	}

	pagination := &mediaprocessingrepository.Pagination{
		Page:  input.Page,
		Limit: input.Limit,
		Sort:  input.Sort,
		Order: input.Order,
	}

	assets, err := s.repo.ListAssets(ctx, tx, uint64(input.ListingIdentityID), filter, pagination)
	if err != nil {
		return dto.ListMediaOutput{}, derrors.Infra("failed to list assets", err)
	}

	count, err := s.repo.CountAssets(ctx, tx, uint64(input.ListingIdentityID), filter)
	if err != nil {
		return dto.ListMediaOutput{}, derrors.Infra("failed to count assets", err)
	}

	return dto.ListMediaOutput{
		Assets:     assets,
		TotalCount: count,
		Page:       input.Page,
		Limit:      input.Limit,
	}, nil
}
