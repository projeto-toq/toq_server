package mysqlproposaladapter

import (
	"context"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/converters"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateStatus updates proposal status-related fields (status, notes, timestamps).
func (a *ProposalAdapter) UpdateStatus(ctx context.Context, proposal proposalmodel.ProposalInterface) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity := converters.ToProposalEntity(proposal)
	now := time.Now().UTC()
	entity.UpdatedAt = now
	proposal.SetUpdatedAt(now)

	query := `UPDATE proposals SET
		status = ?,
		rejection_reason = ?,
		owner_notes = ?,
		accepted_at = ?,
		rejected_at = ?,
		cancelled_at = ?,
		updated_at = ?
	WHERE id = ? AND deleted = FALSE`

	_, execErr := a.ExecContext(ctx, nil, "update_proposal_status", query,
		entity.Status,
		entity.RejectionReason,
		entity.OwnerNotes,
		entity.AcceptedAt,
		entity.RejectedAt,
		entity.CancelledAt,
		entity.UpdatedAt,
		entity.ID,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.proposal.update_status.exec_error", "proposal_id", entity.ID, "err", execErr)
		return fmt.Errorf("update proposal status: %w", execErr)
	}

	return nil
}
