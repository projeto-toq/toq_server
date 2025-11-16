package mediaprocessingconverters

import (
	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// BatchDomainToEntity converte o domínio em entidade SQL.
func BatchDomainToEntity(batch mediaprocessingmodel.MediaBatch) mediaprocessingentities.BatchEntity {
	metadata := batch.StatusMetadata()
	return mediaprocessingentities.BatchEntity{
		ID:              batch.ID(),
		ListingID:       batch.ListingID(),
		Reference:       batch.Reference(),
		Status:          batch.Status().String(),
		StatusMessage:   nullString(metadata.Message),
		StatusReason:    nullString(metadata.Reason),
		StatusDetails:   encodeStringMap(metadata.Details),
		StatusUpdatedBy: metadata.UpdatedBy,
		StatusUpdatedAt: metadata.UpdatedAt,
		DeletedAt:       nullTimeFromPtr(batch.DeletedAt()),
	}
}
