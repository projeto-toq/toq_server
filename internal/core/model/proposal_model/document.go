package proposalmodel

import "time"

// ProposalDocumentInterface now represents the stored PDF metadata and bytes.
type ProposalDocumentInterface interface {
	ID() int64
	SetID(int64)
	ProposalID() int64
	SetProposalID(int64)
	FileName() string
	SetFileName(string)
	MimeType() string
	SetMimeType(string)
	FileSizeBytes() int64
	SetFileSizeBytes(int64)
	FileData() []byte
	SetFileData([]byte)
	UploadedAt() time.Time
	SetUploadedAt(time.Time)
}

type proposalDocument struct {
	id            int64
	proposalID    int64
	fileName      string
	mimeType      string
	fileSizeBytes int64
	fileData      []byte
	uploadedAt    time.Time
}

// NewProposalDocument instantiates a domain document entity.
func NewProposalDocument() ProposalDocumentInterface {
	return &proposalDocument{}
}

func (d *proposalDocument) ID() int64 { return d.id }
func (d *proposalDocument) SetID(id int64) {
	d.id = id
}
func (d *proposalDocument) ProposalID() int64 { return d.proposalID }
func (d *proposalDocument) SetProposalID(id int64) {
	d.proposalID = id
}
func (d *proposalDocument) FileName() string { return d.fileName }
func (d *proposalDocument) SetFileName(name string) {
	d.fileName = name
}
func (d *proposalDocument) MimeType() string { return d.mimeType }
func (d *proposalDocument) SetMimeType(mime string) {
	d.mimeType = mime
}
func (d *proposalDocument) FileSizeBytes() int64 { return d.fileSizeBytes }
func (d *proposalDocument) SetFileSizeBytes(size int64) {
	d.fileSizeBytes = size
}
func (d *proposalDocument) FileData() []byte { return d.fileData }
func (d *proposalDocument) SetFileData(data []byte) {
	if data == nil {
		d.fileData = nil
		return
	}
	buffer := make([]byte, len(data))
	copy(buffer, data)
	d.fileData = buffer
}
func (d *proposalDocument) UploadedAt() time.Time { return d.uploadedAt }
func (d *proposalDocument) SetUploadedAt(ts time.Time) {
	d.uploadedAt = ts
}
