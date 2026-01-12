package mysqlproposaladapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/entities"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *ProposalAdapter) fetchProposal(ctx context.Context, tx *sql.Tx, proposalID int64, forUpdate bool) (proposalmodel.ProposalInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT
        p.id,
        p.listing_identity_id,
        p.realtor_id,
        p.owner_id,
        p.proposal_text,
        p.rejection_reason,
        p.status,
        p.accepted_at,
        p.rejected_at,
	p.cancelled_at,
	p.first_owner_action_at,
	p.created_at,
        p.deleted,
        (
            SELECT COUNT(1)
            FROM proposal_documents d
            WHERE d.proposal_id = p.id
        ) AS documents_count
    FROM proposals p
    WHERE p.id = ? AND p.deleted = 0`

	if forUpdate {
		query += " FOR UPDATE"
	}

	row := a.QueryRowContext(ctx, tx, "get_proposal_by_id", query, proposalID)
	entity := entities.ProposalEntity{}
	if scanErr := row.Scan(
		&entity.ID,
		&entity.ListingIdentityID,
		&entity.RealtorID,
		&entity.OwnerID,
		&entity.ProposalText,
		&entity.RejectionReason,
		&entity.Status,
		&entity.AcceptedAt,
		&entity.RejectedAt,
		&entity.CancelledAt,
		&entity.FirstOwnerAction,
		&entity.CreatedAt,
		&entity.Deleted,
		&entity.DocumentsCount,
	); scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.proposal.get.scan_error", "proposal_id", proposalID, "err", scanErr)
		return nil, fmt.Errorf("get proposal: %w", scanErr)
	}

	return converters.ToProposalModel(entity), nil
}
