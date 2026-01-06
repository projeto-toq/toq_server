package mysqlproposaladapter

import (
	"context"
	"database/sql"

	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
)

// GetProposalByID fetches a proposal row without acquiring locks.
func (a *ProposalAdapter) GetProposalByID(ctx context.Context, tx *sql.Tx, proposalID int64) (proposalmodel.ProposalInterface, error) {
	return a.fetchProposal(ctx, tx, proposalID, false)
}
