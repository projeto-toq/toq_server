package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateProposalFlags toggles proposal status flags on listing identities within the current transaction.
func (la *ListingAdapter) UpdateProposalFlags(ctx context.Context, tx *sql.Tx, input listingrepository.ProposalFlagsUpdate) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE listing_identities
		SET has_pending_proposal = ?,
			has_accepted_proposal = ?,
			accepted_proposal_id = ?
		WHERE id = ? AND deleted = 0`

	defer la.ObserveOnComplete("update", query)()

	pending := 0
	if input.HasPending {
		pending = 1
	}
	accepted := 0
	if input.HasAccepted {
		accepted = 1
	}

	result, execErr := tx.ExecContext(ctx, query, pending, accepted, input.AcceptedProposalID, input.ListingIdentityID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.listing.update_proposal_flags.exec_error", "err", execErr, "listing_identity_id", input.ListingIdentityID)
		return fmt.Errorf("update listing proposal flags: %w", execErr)
	}

	affected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.listing.update_proposal_flags.rows_error", "err", rowsErr, "listing_identity_id", input.ListingIdentityID)
		return fmt.Errorf("rows affected for proposal flags: %w", rowsErr)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
