package mysqllistingfavoriteadapter

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// Add links a user to a listing identity as favorite. Idempotent via ON DUPLICATE KEY.
func (a *ListingFavoriteAdapter) Add(ctx context.Context, tx *sql.Tx, userID, listingIdentityID int64) error {
	ctx, spanEnd, _ := utils.GenerateTracer(ctx)
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `INSERT INTO listing_favorites (user_id, listing_identity_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE listing_identity_id = listing_identity_id`

	if _, err := a.ExecContext(ctx, tx, "insert", query, userID, listingIdentityID); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.listing_favorite.add.exec_error", "user_id", userID, "listing_identity_id", listingIdentityID, "err", err)
		return err
	}

	slog.Debug("mysql.listing_favorite.add.ok", "user_id", userID, "listing_identity_id", listingIdentityID)
	return nil
}
