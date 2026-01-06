package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetListingIdentityByID(ctx context.Context, tx *sql.Tx, identityID int64) (listingrepository.ListingIdentityRecord, error) {
	record := listingrepository.ListingIdentityRecord{}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return record, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id,
		listing_uuid,
		user_id,
		code,
		active_version_id,
		deleted,
		has_pending_proposal,
		has_accepted_proposal,
		accepted_proposal_id
	FROM listing_identities
	WHERE id = ?`

	var (
		activeVersion sql.NullInt64
		deleted       sql.NullInt16
		codeValue     sql.NullInt64
		hasPending    sql.NullInt16
		hasAccepted   sql.NullInt16
		acceptedID    sql.NullInt64
	)

	row := la.QueryRowContext(ctx, tx, "select", query, identityID)
	if scanErr := row.Scan(
		&record.ID,
		&record.UUID,
		&record.UserID,
		&codeValue,
		&activeVersion,
		&deleted,
		&hasPending,
		&hasAccepted,
		&acceptedID,
	); scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return record, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.listing.get_listing_identity_by_id.scan_error", "error", scanErr, "identity_id", identityID)
		return record, fmt.Errorf("scan listing identity by id: %w", scanErr)
	}

	record.ActiveVersionID = activeVersion
	if codeValue.Valid {
		record.Code = uint32(codeValue.Int64)
	}
	if deleted.Valid {
		record.Deleted = deleted.Int16 == 1
	}
	record.HasPendingProposal = hasPending.Valid && hasPending.Int16 == 1
	record.HasAcceptedProposal = hasAccepted.Valid && hasAccepted.Int16 == 1
	record.AcceptedProposalID = acceptedID

	return record, nil
}
