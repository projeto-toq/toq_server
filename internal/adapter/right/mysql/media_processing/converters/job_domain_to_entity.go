package mediaprocessingconverters

import (
	"database/sql"
	"encoding/json"

	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// JobDomainToEntity converte o domínio em entidade listing_media_jobs.
func JobDomainToEntity(job mediaprocessingmodel.MediaProcessingJob) mediaprocessingentities.JobEntity {
	return mediaprocessingentities.JobEntity{
		ID:          job.ID(),
		BatchID:     job.BatchID(),
		ListingID:   job.ListingID(),
		Status:      string(job.Status()),
		Provider:    string(job.Provider()),
		ExternalID:  nullString(job.ExternalID()),
		Payload:     encodeJobPayload(job.Payload()),
		RetryCount:  job.RetryCount(),
		StartedAt:   nullTimeFromPtr(job.StartedAt()),
		CompletedAt: nullTimeFromPtr(job.CompletedAt()),
		LastError:   nullString(job.LastError()),
		CallbackRaw: nullString(job.CallbackBody()),
	}
}

func encodeJobPayload(payload mediaprocessingmodel.MediaProcessingJobPayload) sql.NullString {
	if isEmptyJobPayload(payload) {
		return sql.NullString{}
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return sql.NullString{}
	}
	return sql.NullString{String: string(bytes), Valid: true}
}

// EncodeJobPayload expõe a serialização de payload.
func EncodeJobPayload(payload mediaprocessingmodel.MediaProcessingJobPayload) sql.NullString {
	return encodeJobPayload(payload)
}

func isEmptyJobPayload(payload mediaprocessingmodel.MediaProcessingJobPayload) bool {
	return payload.RawKey == "" &&
		payload.ProcessedKey == "" &&
		payload.ThumbnailKey == "" &&
		len(payload.Outputs) == 0 &&
		payload.ErrorCode == "" &&
		payload.ErrorMessage == ""
}
