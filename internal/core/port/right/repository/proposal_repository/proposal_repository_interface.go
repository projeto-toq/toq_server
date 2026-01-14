package proposalrepository

import (
	"context"
	"database/sql"

	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
)

// Repository exposes CRUD operations for proposals and documents (Seção 7.3 do guia).
type Repository interface {
	CreateProposal(ctx context.Context, tx *sql.Tx, proposal proposalmodel.ProposalInterface) error
	UpdateProposalText(ctx context.Context, tx *sql.Tx, proposal proposalmodel.ProposalInterface) error
	UpdateProposalStatus(ctx context.Context, tx *sql.Tx, proposal proposalmodel.ProposalInterface, expected proposalmodel.Status) error
	GetProposalByID(ctx context.Context, tx *sql.Tx, proposalID int64) (proposalmodel.ProposalInterface, error)
	GetProposalByIDForUpdate(ctx context.Context, tx *sql.Tx, proposalID int64) (proposalmodel.ProposalInterface, error)
	ListProposals(ctx context.Context, tx *sql.Tx, filter proposalmodel.ListFilter) (proposalmodel.ListResult, error)
	CreateDocument(ctx context.Context, tx *sql.Tx, document proposalmodel.ProposalDocumentInterface) error
	ListDocuments(ctx context.Context, tx *sql.Tx, proposalID int64, includeBlob bool) ([]proposalmodel.ProposalDocumentInterface, error)
	ListDocumentsByProposalIDs(ctx context.Context, tx *sql.Tx, proposalIDs []int64, includeBlob bool) (map[int64][]proposalmodel.ProposalDocumentInterface, error)
	ListRealtorSummaries(ctx context.Context, tx *sql.Tx, realtorIDs []int64) ([]proposalmodel.RealtorSummary, error)
	ListOwnerSummaries(ctx context.Context, tx *sql.Tx, ownerIDs []int64) ([]proposalmodel.OwnerSummary, error)
}
