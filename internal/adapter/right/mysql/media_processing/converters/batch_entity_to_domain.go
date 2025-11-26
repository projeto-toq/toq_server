package mediaprocessingconverters

import (
	"encoding/json"
	"time"

	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// BatchEntityToDomain converte uma entidade SQL para o domínio.
func BatchEntityToDomain(entity mediaprocessingentities.BatchEntity) (mediaprocessingmodel.MediaBatch, error) {
	// Deserializa manifesto para obter reference
	var manifest BatchManifest
	if len(entity.UploadManifestJSON) > 0 {
		if err := json.Unmarshal(entity.UploadManifestJSON, &manifest); err != nil {
			return mediaprocessingmodel.MediaBatch{}, err
		}
	}

	// Deserializa detalhes
	details := make(map[string]string)
	if len(entity.ProcessingMetadataJSON) > 0 {
		if err := json.Unmarshal(entity.ProcessingMetadataJSON, &details); err != nil {
			return mediaprocessingmodel.MediaBatch{}, err
		}
	}

	metadata := mediaprocessingmodel.BatchStatusMetadata{
		Message:   entity.ErrorDetail.String,
		Reason:    entity.ErrorCode.String,
		Details:   details,
		UpdatedBy: entity.PhotographerUserID,
		// Usa ProcessingFinishedAt ou StartedAt ou ReceivedAt como fallback para UpdatedAt
		UpdatedAt: resolveUpdatedAt(entity),
	}

	record := mediaprocessingmodel.MediaBatchRecord{
		ID:             entity.ID,
		ListingID:      entity.ListingID,
		Reference:      manifest.BatchReference,
		Status:         mediaprocessingmodel.BatchStatus(entity.Status),
		StatusMetadata: metadata,
		Assets:         []mediaprocessingmodel.MediaAsset{}, // Assets carregados separadamente
		DeletedAt:      timePtrFromNull(entity.DeletedAt),
	}

	return mediaprocessingmodel.RestoreMediaBatch(record), nil
}

func resolveUpdatedAt(e mediaprocessingentities.BatchEntity) time.Time {
	if e.ProcessingFinishedAt.Valid {
		return e.ProcessingFinishedAt.Time
	}
	if e.ProcessingStartedAt.Valid {
		return e.ProcessingStartedAt.Time
	}
	if e.ReceivedAt.Valid {
		return e.ReceivedAt.Time
	}
	return time.Time{} // Fallback
}
