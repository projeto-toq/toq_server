package mediaprocessingconverters

import (
	"encoding/json"

	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// BatchDomainToEntity converte o domínio em entidade SQL.
func BatchDomainToEntity(batch mediaprocessingmodel.MediaBatch) (mediaprocessingentities.BatchEntity, error) {
	metadata := batch.StatusMetadata()

	// Serializa o manifesto para extrair o reference
	manifest := BatchManifest{
		BatchReference: batch.Reference(),
	}
	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		return mediaprocessingentities.BatchEntity{}, err
	}

	// Serializa detalhes de processamento
	var processingMeta []byte
	if len(metadata.Details) > 0 {
		processingMeta, err = json.Marshal(metadata.Details)
		if err != nil {
			return mediaprocessingentities.BatchEntity{}, err
		}
	}

	return mediaprocessingentities.BatchEntity{
		ID:                     batch.ID(),
		ListingID:              batch.ListingID(),
		PhotographerUserID:     metadata.UpdatedBy, // Assume criador como fotógrafo
		Status:                 batch.Status().String(),
		UploadManifestJSON:     manifestBytes,
		ProcessingMetadataJSON: processingMeta,
		ErrorCode:              nullString(metadata.Reason),
		ErrorDetail:            nullString(metadata.Message),
		// Timestamps de processamento podem ser inferidos ou deixados NULL dependendo do status
		DeletedAt: nullTimeFromPtr(batch.DeletedAt()),
	}, nil
}
