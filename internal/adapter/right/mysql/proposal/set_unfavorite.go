package mysqlproposaladapter

import "context"

// SetUnfavorite clears favorite flag on a proposal.
func (a *ProposalAdapter) SetUnfavorite(ctx context.Context, proposalID int64) error {
	return a.SetFavorite(ctx, proposalID, false)
}
