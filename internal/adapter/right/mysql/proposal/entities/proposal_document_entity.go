package entities

import "time"

// ProposalDocumentEntity mirrors the proposal_documents table including the blob payload.
type ProposalDocumentEntity struct {
	ID            int64
	ProposalID    int64
	FileName      string
	MimeType      string
	FileSizeBytes int64
	FileBlob      []byte
	UploadedAt    time.Time
}
