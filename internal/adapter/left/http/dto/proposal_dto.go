package dto

import "time"

// CreateProposalRequest represents the realtor payload to submit a new proposal.
type CreateProposalRequest struct {
	ListingIdentityID int64                   `json:"listingIdentityId" binding:"required,min=1" example:"981"`
	ProposalText      string                  `json:"proposalText" binding:"required,min=1,max=5000" example:"Gostaria de propor pagamento em 30 dias"`
	Document          *ProposalDocumentUpload `json:"document,omitempty"`
}

// UpdateProposalRequest allows editing a pending proposal text/document.
type UpdateProposalRequest struct {
	ProposalID   int64                   `json:"proposalId" binding:"required,min=1" example:"120"`
	ProposalText string                  `json:"proposalText" binding:"required,min=1,max=5000"`
	Document     *ProposalDocumentUpload `json:"document,omitempty"`
}

// ProposalDocumentUpload carries the base64 PDF metadata limited to 1MB.
type ProposalDocumentUpload struct {
	FileName      string `json:"fileName" binding:"required,min=1,max=120" example:"proposta.pdf"`
	MimeType      string `json:"mimeType" binding:"required,oneof=application/pdf" example:"application/pdf"`
	Base64Payload string `json:"base64Payload" binding:"required"`
}

// CancelProposalRequest is used by realtors before owner acceptance.
type CancelProposalRequest struct {
	ProposalID int64 `json:"proposalId" binding:"required,min=1"`
}

// AcceptProposalRequest is triggered by the owner.
type AcceptProposalRequest struct {
	ProposalID int64 `json:"proposalId" binding:"required,min=1"`
}

// RejectProposalRequest stores the owner reason for refusal.
type RejectProposalRequest struct {
	ProposalID int64  `json:"proposalId" binding:"required,min=1"`
	Reason     string `json:"reason" binding:"required,min=1,max=500"`
}

// ListProposalsQuery is shared by realtor/owner GET endpoints.
type ListProposalsQuery struct {
	Statuses          []string `form:"status" binding:"omitempty,dive,oneof=pending accepted refused cancelled"`
	ListingIdentityID int64    `form:"listingIdentityId" binding:"omitempty,min=1"`
	Page              int      `form:"page" binding:"omitempty,min=1" default:"1"`
	PageSize          int      `form:"pageSize" binding:"omitempty,min=1,max=100" default:"20"`
}

// GetProposalDetailRequest returns the full payload including documents.
type GetProposalDetailRequest struct {
	ProposalID int64 `json:"proposalId" binding:"required,min=1"`
}

// ProposalResponse summarizes proposal information for list views, including realtor enrichment and timeline metadata.
type ProposalResponse struct {
	ID                int64              `json:"id"`
	ListingIdentityID int64              `json:"listingIdentityId"`
	Listing           *ListingSummaryDTO `json:"listing,omitempty"`
	Status            string             `json:"status"`
	ProposalText      string             `json:"proposalText"`
	RejectionReason   *string            `json:"rejectionReason,omitempty"`
	AcceptedAt        *time.Time         `json:"acceptedAt,omitempty"`
	RejectedAt        *time.Time         `json:"rejectedAt,omitempty"`
	CancelledAt       *time.Time         `json:"cancelledAt,omitempty"`
	OwnerViewed       bool               `json:"ownerViewed"`
	OwnerViewedAt     *time.Time         `json:"ownerViewedAt,omitempty"`
	// CreatedAt reflects when the realtor submitted the proposal.
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	// ReceivedAt mirrors the first owner action timestamp, signaling that the owner viewed the proposal.
	ReceivedAt *time.Time `json:"receivedAt,omitempty"`
	// RespondedAt is the earliest timestamp among accepted/rejected/cancelled transitions.
	RespondedAt    *time.Time                 `json:"respondedAt,omitempty"`
	DocumentsCount int                        `json:"documentsCount"`
	Documents      []ProposalDocumentResponse `json:"documents"`
	Realtor        ProposalRealtorResponse    `json:"realtor"`
	Owner          ProposalOwnerResponse      `json:"owner"`
}

// ProposalDocumentResponse exposes metadata and optional base64 payload.
type ProposalDocumentResponse struct {
	ID            int64  `json:"id"`
	FileName      string `json:"fileName"`
	MimeType      string `json:"mimeType"`
	FileSizeBytes int64  `json:"fileSizeBytes"`
	Base64Payload string `json:"base64Payload,omitempty"`
}

// ProposalRealtorResponse describes enriched realtor metadata exposed to owners and realtors.
type ProposalRealtorResponse struct {
	Name     string `json:"name"`
	Nickname string `json:"nickname,omitempty"`
	// AccountAgeMonths represents how long (in months) the realtor has been active in TOQ.
	AccountAgeMonths int `json:"accountAgeMonths"`
	// ProposalsCreated is the lifetime counter of proposals authored.
	ProposalsCreated int64 `json:"proposalsCreated"`
	// AcceptedProposals tracks how many proposals from this realtor owners accepted.
	AcceptedProposals int64 `json:"acceptedProposals"`
	// PhotoURL is a signed download URL pointing to the realtor avatar variant consumed by the clients.
	PhotoURL string `json:"photoUrl,omitempty"`
}

// ProposalOwnerResponse exposes owner profile metadata and engagement metrics.
type ProposalOwnerResponse struct {
	ID                 int64  `json:"id"`
	FullName           string `json:"fullName"`
	MemberSinceMonths  int    `json:"memberSinceMonths"`
	PhotoURL           string `json:"photoUrl,omitempty"`
	ProposalAvgSeconds *int64 `json:"proposalAverageSeconds,omitempty"`
	VisitAvgSeconds    *int64 `json:"visitAverageSeconds,omitempty"`
}

// ProposalDetailResponse aggregates summary + documents and owner metadata.
type ProposalDetailResponse struct {
	Proposal  ProposalResponse           `json:"proposal"`
	Documents []ProposalDocumentResponse `json:"documents"`
	Realtor   ProposalRealtorResponse    `json:"realtor"`
	Owner     ProposalOwnerResponse      `json:"owner"`
}

// ListProposalsResponse is returned by both realtor/owner endpoints.
type ListProposalsResponse struct {
	Items []ProposalResponse `json:"items"`
	Total int64              `json:"total"`
}
