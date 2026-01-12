package mysqlproposaladapter

import (
	"context"
	"database/sql"
	"fmt"

	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateProposalStatus updates the status columns enforcing optimistic locking via expected status.
func (a *ProposalAdapter) UpdateProposalStatus(ctx context.Context, tx *sql.Tx, proposal proposalmodel.ProposalInterface, expected proposalmodel.Status) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE proposals
	        SET status = ?, rejection_reason = ?, accepted_at = ?, rejected_at = ?, cancelled_at = ?, first_owner_action_at = ?
	    WHERE id = ? AND status = ? AND deleted = 0`

	result, execErr := a.ExecContext(ctx, tx, "update_proposal_status", query,
		proposal.Status(),
		proposal.RejectionReason(),
		proposal.AcceptedAt(),
		proposal.RejectedAt(),
		proposal.CancelledAt(),
		proposal.FirstOwnerActionAt(),
		proposal.ID(),
		expected,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.proposal.update_status.exec_error", "proposal_id", proposal.ID(), "err", execErr)
		return fmt.Errorf("update proposal status: %w", execErr)
	}

	rows, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.proposal.update_status.rows_error", "proposal_id", proposal.ID(), "err", rowsErr)
		return fmt.Errorf("update proposal status rows: %w", rowsErr)
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
