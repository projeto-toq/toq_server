package mediaprocessingconverters

import (
	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// AssetEntityToDomain converte uma entidade SQL em domínio.
func AssetEntityToDomain(entity mediaprocessingentities.AssetEntity) mediaprocessingmodel.MediaAsset {
	record := mediaprocessingmodel.MediaAssetRecord{
		ID:             entity.ID,
		BatchID:        entity.BatchID,
		ListingID:      entity.ListingID,
		AssetType:      mediaprocessingmodel.MediaAssetType(entity.AssetType),
		Orientation:    mediaprocessingmodel.MediaAssetOrientation(entity.Orientation.String),
		Filename:       entity.Filename,
		ContentType:    entity.ContentType,
		Sequence:       entity.Sequence,
		SizeInBytes:    entity.SizeInBytes,
		Checksum:       entity.Checksum,
		RawObjectKey:   entity.RawObjectKey,
		ProcessedKey:   entity.ProcessedKey.String,
		ThumbnailKey:   entity.ThumbnailKey.String,
		Width:          uint16(entity.Width.Int64),
		Height:         uint16(entity.Height.Int64),
		DurationMillis: uint32(entity.DurationMillis.Int64),
		Metadata:       decodeStringMap(entity.Metadata),
	}

	return mediaprocessingmodel.RestoreMediaAsset(record)
}
