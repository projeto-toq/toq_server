package mediaprocessingservice

import (
	"context"
	"database/sql"
	"errors"

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

	zipBundle, err := s.latestZipBundle(ctx, tx, uint64(input.ListingIdentityID))
	if err != nil {
		return dto.ListMediaOutput{}, err
	}

	return dto.ListMediaOutput{
		Assets:     assets,
		TotalCount: count,
		Page:       input.Page,
		Limit:      input.Limit,
		ZipBundle:  zipBundle,
	}, nil
}

func (s *mediaProcessingService) latestZipBundle(ctx context.Context, tx *sql.Tx, listingIdentityID uint64) (*dto.ListMediaZipBundle, error) {
	job, err := s.repo.GetLatestFinalizationJob(ctx, tx, listingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, derrors.Infra("failed to fetch finalization job", err)
	}

	payload := job.Payload()
	bundleKey := ""
	if len(payload.ZipBundles) > 0 {
		bundleKey = payload.ZipBundles[0]
	}

	if bundleKey == "" && payload.ZipSizeBytes == 0 && payload.UnzippedSizeBytes == 0 && payload.AssetsZipped == 0 {
		return nil, nil
	}

	return &dto.ListMediaZipBundle{
		BundleKey:               bundleKey,
		AssetsCount:             payload.AssetsZipped,
		ZipSizeBytes:            payload.ZipSizeBytes,
		EstimatedExtractedBytes: payload.UnzippedSizeBytes,
		CompletedAt:             job.CompletedAt(),
	}, nil
}
