package mysqlownermetricsadapter

import (
	"context"
	"database/sql"
	"fmt"

	ownermetricsrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/owner_metrics_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpsertVisitResponse aggregates visit SLA metrics per owner.
func (a *OwnerMetricsAdapter) UpsertVisitResponse(ctx context.Context, tx *sql.Tx, input ownermetricsrepository.VisitResponseInput) error {
	if input.OwnerID <= 0 {
		return fmt.Errorf("owner id must be positive")
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `INSERT INTO owner_response_metrics (
		user_id,
		visit_avg_response_time_seconds,
		visit_total_responses,
		visit_last_response_at
	) VALUES (?, ?, 1, ?)
	ON DUPLICATE KEY UPDATE
		visit_avg_response_time_seconds = CASE
			WHEN owner_response_metrics.visit_avg_response_time_seconds IS NULL OR owner_response_metrics.visit_total_responses = 0 THEN VALUES(visit_avg_response_time_seconds)
			ELSE FLOOR((owner_response_metrics.visit_avg_response_time_seconds * owner_response_metrics.visit_total_responses + VALUES(visit_avg_response_time_seconds)) / (owner_response_metrics.visit_total_responses + 1))
		END,
		visit_total_responses = owner_response_metrics.visit_total_responses + 1,
		visit_last_response_at = GREATEST(COALESCE(owner_response_metrics.visit_last_response_at, VALUES(visit_last_response_at)), VALUES(visit_last_response_at))`

	defer a.ObserveOnComplete("upsert_visit_response", query)()

	if _, execErr := a.ExecContext(ctx, tx, "upsert_visit_response", query,
		input.OwnerID,
		input.DeltaSeconds,
		input.RespondedAt,
	); execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.owner_metrics.visit.upsert_error", "owner_id", input.OwnerID, "err", execErr)
		return fmt.Errorf("upsert owner visit metrics: %w", execErr)
	}

	return nil
}
