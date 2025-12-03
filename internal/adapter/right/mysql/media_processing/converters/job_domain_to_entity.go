package mediaprocessingconverters

import (
	"database/sql"
	"encoding/json"

	mediaprocessingentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/media_processing/entities"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// JobDomainToEntity converts domain to DB entity.
func JobDomainToEntity(job mediaprocessingmodel.MediaProcessingJob) mediaprocessingentities.JobEntity {
	return mediaprocessingentities.JobEntity{
		ID:                job.ID(),
		ListingIdentityID: job.ListingIdentityID(),
		Status:            string(job.Status()),
		Provider:          string(job.Provider()),
		ExternalID:        nullString(job.ExternalID()),
		Payload:           EncodeJobPayload(job.Payload()),
		StartedAt:         nullTimeFromPtr(job.StartedAt()),
		FinishedAt:        nullTimeFromPtr(job.CompletedAt()),
		LastError:         nullString(job.LastError()),
		CallbackBody:      nullString(job.CallbackBody()),
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
		payload.ErrorMessage == "" &&
		len(payload.ZipBundles) == 0 &&
		payload.AssetsZipped == 0
}
