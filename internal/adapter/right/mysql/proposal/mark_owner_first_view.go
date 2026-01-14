package mysqlproposaladapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// MarkOwnerFirstView sets the first_owner_action_at timestamp when an owner opens a proposal for the first time.
// It is idempotent: if the column is already populated, no update is performed.
func (a *ProposalAdapter) MarkOwnerFirstView(ctx context.Context, tx *sql.Tx, proposalID int64, ownerID int64, seenAt time.Time) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE proposals
		SET first_owner_action_at = ?
		WHERE id = ? AND owner_id = ? AND deleted = 0 AND first_owner_action_at IS NULL`

	result, execErr := a.ExecContext(ctx, tx, "mark_owner_first_view", query, seenAt, proposalID, ownerID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.proposal.mark_first_view.exec_error", "proposal_id", proposalID, "owner_id", ownerID, "err", execErr)
		return fmt.Errorf("mark owner first view: %w", execErr)
	}

	if _, rowsErr := result.RowsAffected(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.proposal.mark_first_view.rows_error", "proposal_id", proposalID, "owner_id", ownerID, "err", rowsErr)
		return fmt.Errorf("mark owner first view rows: %w", rowsErr)
	}

	return nil
}
