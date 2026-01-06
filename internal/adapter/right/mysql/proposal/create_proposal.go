package mysqlproposaladapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/converters"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateProposal inserts a new proposal row using the provided transaction and updates domain timestamps/ID.
func (a *ProposalAdapter) CreateProposal(ctx context.Context, tx *sql.Tx, proposal proposalmodel.ProposalInterface) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity := converters.ToProposalEntity(proposal)

	query := `INSERT INTO proposals (
		listing_identity_id,
		realtor_id,
		owner_id,
		status,
		proposal_text,
		rejection_reason,
		accepted_at,
		rejected_at,
		cancelled_at,
		deleted
	) VALUES (?,?,?,?,?,?,?,?,?,0)`

	result, execErr := a.ExecContext(ctx, tx, "insert_proposal", query,
		entity.ListingIdentityID,
		entity.RealtorID,
		entity.OwnerID,
		entity.Status,
		entity.ProposalText,
		entity.RejectionReason,
		entity.AcceptedAt,
		entity.RejectedAt,
		entity.CancelledAt,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.proposal.create.exec_error", "listing_identity_id", entity.ListingIdentityID, "err", execErr)
		return fmt.Errorf("create proposal: %w", execErr)
	}

	id, idErr := result.LastInsertId()
	if idErr != nil {
		utils.SetSpanError(ctx, idErr)
		logger.Error("mysql.proposal.create.last_insert_id_error", "err", idErr)
		return fmt.Errorf("proposal last insert id: %w", idErr)
	}

	proposal.SetID(id)
	return nil
}
