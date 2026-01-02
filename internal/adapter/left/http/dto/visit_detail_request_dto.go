package dto

// GetVisitDetailRequest carries the identifier to fetch a single visit detail.
type GetVisitDetailRequest struct {
	// VisitID is the unique identifier of the visit to be retrieved.
	VisitID int64 `json:"visitId" binding:"required,gt=0" example:"456"`
}
