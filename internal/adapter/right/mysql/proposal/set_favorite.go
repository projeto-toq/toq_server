package mysqlproposaladapter

import (
	"context"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// SetFavorite sets or clears favorite flag on a proposal.
func (a *ProposalAdapter) SetFavorite(ctx context.Context, proposalID int64, value bool) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE proposals SET is_favorite = ?, updated_at = NOW() WHERE id = ? AND deleted = FALSE`
	_, execErr := a.ExecContext(ctx, nil, "set_proposal_favorite", query, value, proposalID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.proposal.favorite.exec_error", "proposal_id", proposalID, "err", execErr)
		return fmt.Errorf("set proposal favorite: %w", execErr)
	}

	return nil
}
