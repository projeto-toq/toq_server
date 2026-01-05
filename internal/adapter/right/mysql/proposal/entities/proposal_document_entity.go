package entities

import "time"

// ProposalDocumentEntity mirrors the proposal_documents table.
type ProposalDocumentEntity struct {
	ID            int64
	ProposalID    int64
	FileName      string
	FileType      string
	FileURL       string
	FileSizeBytes int64
	UploadedAt    time.Time
}
