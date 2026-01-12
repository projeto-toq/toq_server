package mysqllistingfavoriteadapter

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// Remove unlinks a user from a listing identity. Idempotent: does not error if nothing is deleted.
func (a *ListingFavoriteAdapter) Remove(ctx context.Context, tx *sql.Tx, userID, listingIdentityID int64) error {
	ctx, spanEnd, _ := utils.GenerateTracer(ctx)
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM listing_favorites WHERE user_id = ? AND listing_identity_id = ?`

	result, err := a.ExecContext(ctx, tx, "delete", query, userID, listingIdentityID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing_favorite.remove.exec_error", "user_id", userID, "listing_identity_id", listingIdentityID, "err", err)
		return err
	}

	if rows, raErr := result.RowsAffected(); raErr != nil {
		utils.SetSpanError(ctx, raErr)
		logger.Error("mysql.listing_favorite.remove.rows_error", "user_id", userID, "listing_identity_id", listingIdentityID, "err", raErr)
		return raErr
	} else {
		slog.Debug("mysql.listing_favorite.remove.rows", "rows", rows, "user_id", userID, "listing_identity_id", listingIdentityID)
	}

	return nil
}
