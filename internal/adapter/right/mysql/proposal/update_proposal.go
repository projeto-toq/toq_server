package mysqlproposaladapter

import (
	"context"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/converters"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateProposal updates mutable proposal fields when still pending.
func (a *ProposalAdapter) UpdateProposal(ctx context.Context, proposal proposalmodel.ProposalInterface) error {
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
		transaction_type = ?,
		payment_method = ?,
		proposed_value = ?,
		original_value = ?,
		down_payment = ?,
		installments = ?,
		accepts_exchange = ?,
		rental_months = ?,
		guarantee_type = ?,
		security_deposit = ?,
		client_name = ?,
		client_phone = ?,
		proposal_notes = ?,
		expires_at = ?,
		updated_at = ?
	WHERE id = ? AND deleted = FALSE`

	_, execErr := a.ExecContext(ctx, nil, "update_proposal", query,
		entity.TransactionType,
		entity.PaymentMethod,
		entity.ProposedValue,
		entity.OriginalValue,
		entity.DownPayment,
		entity.Installments,
		entity.AcceptsExchange,
		entity.RentalMonths,
		entity.GuaranteeType,
		entity.SecurityDeposit,
		entity.ClientName,
		entity.ClientPhone,
		entity.ProposalNotes,
		entity.ExpiresAt,
		entity.UpdatedAt,
		entity.ID,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.proposal.update.exec_error", "proposal_id", entity.ID, "err", execErr)
		return fmt.Errorf("update proposal: %w", execErr)
	}

	return nil
}
