package mysqlproposaladapter

import (
	"context"
	"database/sql"

	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
)

// GetProposalByIDForUpdate fetches and locks a proposal row inside the active transaction.
func (a *ProposalAdapter) GetProposalByIDForUpdate(ctx context.Context, tx *sql.Tx, proposalID int64) (proposalmodel.ProposalInterface, error) {
	return a.fetchProposal(ctx, tx, proposalID, true)
}
