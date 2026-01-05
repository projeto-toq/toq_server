package mysqlproposaladapter

import (
	"context"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/converters"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateProposal persists a new proposal row and updates the domain object with the generated ID.
func (a *ProposalAdapter) CreateProposal(ctx context.Context, proposal proposalmodel.ProposalInterface) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity := converters.ToProposalEntity(proposal)
	now := time.Now().UTC()
	if entity.CreatedAt.IsZero() {
		entity.CreatedAt = now
	}
	entity.UpdatedAt = now

	query := `INSERT INTO proposals (
		listing_identity_id,
		realtor_id,
		owner_id,
		transaction_type,
		payment_method,
		proposed_value,
		original_value,
		down_payment,
		installments,
		accepts_exchange,
		rental_months,
		guarantee_type,
		security_deposit,
		client_name,
		client_phone,
		proposal_notes,
		owner_notes,
		rejection_reason,
		status,
		expires_at,
		accepted_at,
		rejected_at,
		cancelled_at,
		is_favorite,
		created_at,
		updated_at,
		deleted
	) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	result, execErr := a.ExecContext(ctx, nil, "insert_proposal", query,
		entity.ListingIdentityID,
		entity.RealtorID,
		entity.OwnerID,
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
		entity.OwnerNotes,
		entity.RejectionReason,
		entity.Status,
		entity.ExpiresAt,
		entity.AcceptedAt,
		entity.RejectedAt,
		entity.CancelledAt,
		entity.IsFavorite,
		entity.CreatedAt,
		entity.UpdatedAt,
		entity.Deleted,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.proposal.create.exec_error", "listing_identity_id", entity.ListingIdentityID, "err", execErr)
		return fmt.Errorf("create proposal: %w", execErr)
	}

	id, idErr := result.LastInsertId()
	if idErr != nil {
		utils.SetSpanError(ctx, idErr)
		logger.Error("mysql.proposal.create.last_insert_id_error", "listing_identity_id", entity.ListingIdentityID, "err", idErr)
		return fmt.Errorf("proposal last insert id: %w", idErr)
	}

	proposal.SetID(id)
	return nil
}
