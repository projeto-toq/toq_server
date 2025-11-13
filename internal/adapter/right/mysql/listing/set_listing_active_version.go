package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) SetListingActiveVersion(ctx context.Context, tx *sql.Tx, identityID int64, versionID int64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE listing_identities SET active_version_id = ? WHERE id = ?`

	if _, execErr := la.ExecContext(ctx, tx, "update", query, versionID, identityID); execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.listing.set_active_version.exec_error", "error", execErr, "identity_id", identityID, "version_id", versionID)
		return fmt.Errorf("exec set listing active version: %w", execErr)
	}

	return nil
}
