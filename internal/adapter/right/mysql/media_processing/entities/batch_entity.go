package mediaprocessingentities

import (
	"database/sql"
	"encoding/json"
)

// BatchEntity espelha a tabela listing_media_batches.
// Mapeamento estrito 1:1 com o banco de dados.
type BatchEntity struct {
	ID                     uint64
	ListingID              uint64
	PhotographerUserID     uint64
	Status                 string
	UploadManifestJSON     json.RawMessage
	ProcessingMetadataJSON json.RawMessage
	ReceivedAt             sql.NullTime
	ProcessingStartedAt    sql.NullTime
	ProcessingFinishedAt   sql.NullTime
	ErrorCode              sql.NullString
	ErrorDetail            sql.NullString
	DeletedAt              sql.NullTime
}
