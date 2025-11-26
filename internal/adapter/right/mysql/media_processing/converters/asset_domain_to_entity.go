package mediaprocessingconverters

import (
	"database/sql"

	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// AssetDomainToEntity converte o domínio em entidade de banco.
func AssetDomainToEntity(asset mediaprocessingmodel.MediaAsset) mediaprocessingentities.AssetEntity {
	orientation := nullString(string(asset.Orientation()))
	width := sql.NullInt64{}
	if asset.Width() > 0 {
		width = sql.NullInt64{Int64: int64(asset.Width()), Valid: true}
	}
	height := sql.NullInt64{}
	if asset.Height() > 0 {
		height = sql.NullInt64{Int64: int64(asset.Height()), Valid: true}
	}
	duration := sql.NullInt64{}
	if asset.DurationMillis() > 0 {
		duration = sql.NullInt64{Int64: int64(asset.DurationMillis()), Valid: true}
	}

	return mediaprocessingentities.AssetEntity{
		ID:             asset.ID(),
		BatchID:        asset.BatchID(),
		AssetType:      string(asset.AssetType()),
		Orientation:    orientation,
		Filename:       asset.Filename(),
		ContentType:    asset.ContentType(),
		Sequence:       asset.Sequence(),
		SizeBytes:      asset.SizeInBytes(),
		Checksum:       asset.Checksum(),
		RawObjectKey:   asset.RawObjectKey(),
		ProcessedKey:   nullString(asset.ProcessedKey()),
		ThumbnailKey:   nullString(asset.ThumbnailKey()),
		Width:          width,
		Height:         height,
		DurationMillis: duration,
		Metadata:       encodeStringMap(asset.Metadata()),
	}
}
