package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateOwnerResponseStats aggregates owner response metrics for a listing identity.
// deltaSeconds should represent the time between visit creation and the owner's first action.
func (la *ListingAdapter) UpdateOwnerResponseStats(ctx context.Context, tx *sql.Tx, identityID int64, deltaSeconds int64, respondedAt time.Time) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
UPDATE listing_identities
SET
    owner_avg_response_time_seconds = CASE
        WHEN owner_avg_response_time_seconds IS NULL OR owner_total_visits_responded = 0 THEN ?
        ELSE FLOOR((owner_avg_response_time_seconds * owner_total_visits_responded + ?) / (owner_total_visits_responded + 1))
    END,
    owner_total_visits_responded = owner_total_visits_responded + 1,
    owner_last_response_at = ?
WHERE id = ? AND deleted = 0
`
	defer la.ObserveOnComplete("update", query)()

	result, execErr := tx.ExecContext(ctx, query, deltaSeconds, deltaSeconds, respondedAt, identityID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.listing.update_owner_response_stats.exec_error", "error", execErr, "listing_identity_id", identityID)
		return fmt.Errorf("update owner response stats: %w", execErr)
	}

	affected, affErr := result.RowsAffected()
	if affErr != nil {
		utils.SetSpanError(ctx, affErr)
		logger.Error("mysql.listing.update_owner_response_stats.rows_affected_error", "error", affErr, "listing_identity_id", identityID)
		return fmt.Errorf("rows affected for update owner response stats: %w", affErr)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
