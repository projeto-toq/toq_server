package mediaprocessingconverters

import (
	"database/sql"

	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// AssetDomainToEntity converts domain to DB entity.
func AssetDomainToEntity(asset mediaprocessingmodel.MediaAsset) mediaprocessingentities.AssetEntity {
	return mediaprocessingentities.AssetEntity{
		ID:             asset.ID(),
		ListingID:      asset.ListingID(),
		AssetType:      string(asset.AssetType()),
		Sequence:       asset.Sequence(),
		Status:         string(asset.Status()),
		S3KeyRaw:       sql.NullString{String: asset.S3KeyRaw(), Valid: asset.S3KeyRaw() != ""},
		S3KeyProcessed: sql.NullString{String: asset.S3KeyProcessed(), Valid: asset.S3KeyProcessed() != ""},
		Title:          sql.NullString{String: asset.Title(), Valid: asset.Title() != ""},
		Metadata:       sql.NullString{String: asset.Metadata(), Valid: asset.Metadata() != ""},
	}
}
