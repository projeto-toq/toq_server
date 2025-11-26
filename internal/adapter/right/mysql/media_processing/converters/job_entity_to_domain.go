package mediaprocessingconverters

import (
	"database/sql"
	"encoding/json"

	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// JobEntityToDomain converte uma entidade listing_media_jobs para o domínio.
func JobEntityToDomain(entity mediaprocessingentities.JobEntity) mediaprocessingmodel.MediaProcessingJob {
	record := mediaprocessingmodel.MediaProcessingJobRecord{
		ID:           entity.ID,
		BatchID:      entity.BatchID,
		ListingID:    entity.ListingID,
		Status:       mediaprocessingmodel.MediaProcessingJobStatus(entity.Status),
		Provider:     mediaprocessingmodel.MediaProcessingProvider(entity.Provider),
		ExternalID:   entity.ExternalID.String,
		Payload:      decodeJobPayload(entity.Payload),
		RetryCount:   0,
		StartedAt:    timePtrFromNull(entity.StartedAt),
		CompletedAt:  timePtrFromNull(entity.FinishedAt),
		LastError:    "",
		CallbackBody: "",
	}

	return mediaprocessingmodel.RestoreMediaProcessingJob(record)
}

func decodeJobPayload(raw sql.NullString) mediaprocessingmodel.MediaProcessingJobPayload {
	if !raw.Valid || raw.String == "" {
		return mediaprocessingmodel.MediaProcessingJobPayload{}
	}
	var payload mediaprocessingmodel.MediaProcessingJobPayload
	if err := json.Unmarshal([]byte(raw.String), &payload); err != nil {
		return mediaprocessingmodel.MediaProcessingJobPayload{}
	}
	if payload.Outputs == nil {
		payload.Outputs = map[string]string{}
	}
	return payload
}
