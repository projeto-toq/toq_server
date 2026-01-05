package proposal_repository

import (
	context "context"

	proposal_model "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
)

// Repository provides proposal persistence operations.
type Repository interface {
	CreateProposal(ctx context.Context, proposal proposal_model.ProposalInterface) error
	UpdateProposal(ctx context.Context, proposal proposal_model.ProposalInterface) error
	UpdateStatus(ctx context.Context, proposal proposal_model.ProposalInterface) error
	ListProposals(ctx context.Context, filter proposal_model.ListFilter) (proposal_model.ListResult, error)
	GetStats(ctx context.Context, filter proposal_model.StatsFilter) (proposal_model.Stats, error)
	CreateDocument(ctx context.Context, document proposal_model.ProposalDocumentInterface) error
	ListDocuments(ctx context.Context, proposalID int64) ([]proposal_model.ProposalDocumentInterface, error)
	SetFavorite(ctx context.Context, proposalID int64, value bool) error
	SetUnfavorite(ctx context.Context, proposalID int64) error
}
