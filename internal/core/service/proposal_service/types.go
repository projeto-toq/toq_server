package proposalservice

import (
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
)

// Actor stores the authenticated user metadata extracted from middlewares.
type Actor struct {
	UserID   int64
	RoleSlug permissionmodel.RoleSlug
}

// CreateProposalInput aggregates validated data from handlers.
type CreateProposalInput struct {
	ListingIdentityID int64
	RealtorID         int64
	ProposalText      string
	Document          *DocumentPayload
}

// UpdateProposalInput extends CreateProposalInput with an ID.
type UpdateProposalInput struct {
	ProposalID   int64
	EditorID     int64
	ProposalText string
	Document     *DocumentPayload
}

// DocumentPayload carries decoded bytes and metadata.
type DocumentPayload struct {
	FileName  string
	MimeType  string
	Bytes     []byte
	SizeBytes int64
}

// StatusChangeInput is reused by cancel/accept/reject flows.
type StatusChangeInput struct {
	ProposalID int64
	Actor      Actor
	Reason     string
}

// ListFilter stores normalized filters for repository queries.
type ListFilter struct {
	Actor     Actor
	Statuses  []proposalmodel.Status
	ListingID *int64
	Page      int
	PageSize  int
}

// ListItem aggregates the proposal metadata, documents and realtor info for list endpoints.
type ListItem struct {
	Proposal  proposalmodel.ProposalInterface
	Documents []proposalmodel.ProposalDocumentInterface
	Realtor   proposalmodel.RealtorSummary
	Listing   listingmodel.ListingInterface
}

// ListResult is returned to handlers before DTO serialization.
type ListResult struct {
	Items []ListItem
	Total int64
}

// DetailInput ensures only owners or authors can inspect documents.
type DetailInput struct {
	ProposalID int64
	Actor      Actor
}

// DetailResult stores proposal and documents.
type DetailResult struct {
	Proposal  proposalmodel.ProposalInterface
	Documents []proposalmodel.ProposalDocumentInterface
	Realtor   proposalmodel.RealtorSummary
	Listing   listingmodel.ListingInterface
}
