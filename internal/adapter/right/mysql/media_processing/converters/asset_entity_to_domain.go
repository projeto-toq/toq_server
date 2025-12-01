package mediaprocessingconverters

import (
	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// AssetEntityToDomain converts DB entity to domain.
func AssetEntityToDomain(entity mediaprocessingentities.AssetEntity) mediaprocessingmodel.MediaAsset {
	asset := mediaprocessingmodel.NewMediaAsset(
		entity.ListingID,
		mediaprocessingmodel.MediaAssetType(entity.AssetType),
		entity.Sequence,
	)
	asset.SetID(entity.ID)
	asset.SetStatus(mediaprocessingmodel.MediaAssetStatus(entity.Status))

	if entity.S3KeyRaw.Valid {
		asset.SetS3KeyRaw(entity.S3KeyRaw.String)
	}
	if entity.S3KeyProcessed.Valid {
		asset.SetS3KeyProcessed(entity.S3KeyProcessed.String)
	}
	if entity.Title.Valid {
		asset.SetTitle(entity.Title.String)
	}
	if entity.Metadata.Valid {
		asset.SetMetadata(entity.Metadata.String)
	}

	return asset
}
