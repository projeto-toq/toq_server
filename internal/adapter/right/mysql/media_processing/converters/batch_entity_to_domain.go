package mediaprocessingconverters

import (
	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// BatchEntityToDomain converte uma entidade SQL para o domínio.
func BatchEntityToDomain(entity mediaprocessingentities.BatchEntity) mediaprocessingmodel.MediaBatch {
	metadata := mediaprocessingmodel.BatchStatusMetadata{
		Message:   entity.StatusMessage.String,
		Reason:    entity.StatusReason.String,
		Details:   decodeStringMap(entity.StatusDetails),
		UpdatedBy: entity.StatusUpdatedBy,
		UpdatedAt: entity.StatusUpdatedAt,
	}

	record := mediaprocessingmodel.MediaBatchRecord{
		ID:             entity.ID,
		ListingID:      entity.ListingID,
		Reference:      entity.Reference,
		Status:         mediaprocessingmodel.BatchStatus(entity.Status),
		StatusMetadata: metadata,
		Assets:         []mediaprocessingmodel.MediaAsset{},
		DeletedAt:      timePtrFromNull(entity.DeletedAt),
	}

	return mediaprocessingmodel.RestoreMediaBatch(record)
}
