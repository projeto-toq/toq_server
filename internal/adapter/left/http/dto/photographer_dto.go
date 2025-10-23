package dto

// UpdateSessionStatusRequest defines the payload for updating a session's status.
type UpdateSessionStatusRequest struct {
	Status string `json:"status" binding:"required" example:"ACCEPTED"`
}
