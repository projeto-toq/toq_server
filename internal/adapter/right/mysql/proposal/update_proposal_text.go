package mysqlproposaladapter

import (
	"context"
	"database/sql"
	"fmt"

	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateProposalText updates the proposal free text while the row is pending.
func (a *ProposalAdapter) UpdateProposalText(ctx context.Context, tx *sql.Tx, proposal proposalmodel.ProposalInterface) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE proposals
	        SET proposal_text = ?
	        WHERE id = ? AND status = 'pending' AND deleted = 0`

	result, execErr := a.ExecContext(ctx, tx, "update_proposal_text", query,
		proposal.ProposalText(),
		proposal.ID(),
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.proposal.update_text.exec_error", "proposal_id", proposal.ID(), "err", execErr)
		return fmt.Errorf("update proposal text: %w", execErr)
	}

	rows, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.proposal.update_text.rows_error", "proposal_id", proposal.ID(), "err", rowsErr)
		return fmt.Errorf("update proposal text rows: %w", rowsErr)
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
