package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) CreateListingIdentity(ctx context.Context, tx *sql.Tx, listing listingmodel.ListingInterface) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `INSERT INTO listing_identities (
        listing_uuid,
        user_id,
        code,
        active_version_id,
        deleted
    ) VALUES (?, ?, ?, ?, ?)`

	activeVersionID := sql.NullInt64{}
	if listing.ActiveVersionID() != 0 {
		activeVersionID = sql.NullInt64{Int64: listing.ActiveVersionID(), Valid: true}
	}

	result, execErr := la.ExecContext(ctx, tx, "insert", query,
		listing.UUID(),
		listing.UserID(),
		listing.Code(),
		activeVersionID,
		false,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.listing.create_listing_identity.exec_error", "error", execErr)
		return fmt.Errorf("exec create listing identity: %w", execErr)
	}

	identityID, lastErr := result.LastInsertId()
	if lastErr != nil {
		utils.SetSpanError(ctx, lastErr)
		logger.Error("mysql.listing.create_listing_identity.last_insert_error", "error", lastErr)
		return fmt.Errorf("last insert id for listing identity: %w", lastErr)
	}

	listing.SetIdentityID(identityID)

	return nil
}
