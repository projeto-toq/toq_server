package proposalmodel

import (
	"database/sql"
	"time"
)

// ProposalDocumentInterface defines a document attached to a proposal.
type ProposalDocumentInterface interface {
	ID() int64
	SetID(id int64)
	ProposalID() int64
	SetProposalID(id int64)
	FileName() string
	SetFileName(v string)
	FileType() string
	SetFileType(v string)
	FileURL() string
	SetFileURL(v string)
	FileSizeBytes() int64
	SetFileSizeBytes(v int64)
	UploadedAt() time.Time
	SetUploadedAt(t time.Time)
}

// ProposalDocument is the concrete implementation of ProposalDocumentInterface.
type ProposalDocument struct {
	id            int64
	proposalID    int64
	fileName      string
	fileType      string
	fileURL       string
	fileSizeBytes int64
	uploadedAt    time.Time
}

// NewProposalDocument builds an empty proposal document domain object.
func NewProposalDocument() ProposalDocumentInterface {
	return &ProposalDocument{}
}

func (d *ProposalDocument) ID() int64      { return d.id }
func (d *ProposalDocument) SetID(id int64) { d.id = id }

func (d *ProposalDocument) ProposalID() int64      { return d.proposalID }
func (d *ProposalDocument) SetProposalID(id int64) { d.proposalID = id }

func (d *ProposalDocument) FileName() string     { return d.fileName }
func (d *ProposalDocument) SetFileName(v string) { d.fileName = v }

func (d *ProposalDocument) FileType() string     { return d.fileType }
func (d *ProposalDocument) SetFileType(v string) { d.fileType = v }

func (d *ProposalDocument) FileURL() string     { return d.fileURL }
func (d *ProposalDocument) SetFileURL(v string) { d.fileURL = v }

func (d *ProposalDocument) FileSizeBytes() int64     { return d.fileSizeBytes }
func (d *ProposalDocument) SetFileSizeBytes(v int64) { d.fileSizeBytes = v }

func (d *ProposalDocument) UploadedAt() time.Time     { return d.uploadedAt }
func (d *ProposalDocument) SetUploadedAt(t time.Time) { d.uploadedAt = t }

// ValidateUploadedAt ensures uploadedAt is set when file metadata is present.
func ValidateUploadedAt(value sql.NullTime) time.Time {
	if value.Valid {
		return value.Time
	}
	return time.Time{}
}
