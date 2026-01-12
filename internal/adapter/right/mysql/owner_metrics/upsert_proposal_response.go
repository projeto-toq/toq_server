package mysqlownermetricsadapter

import (
	"context"
	"database/sql"
	"fmt"

	ownermetricsrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/owner_metrics_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpsertProposalResponse aggregates proposal SLA metrics per owner.
func (a *OwnerMetricsAdapter) UpsertProposalResponse(ctx context.Context, tx *sql.Tx, input ownermetricsrepository.ProposalResponseInput) error {
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
		proposal_avg_response_time_seconds,
		proposal_total_responses,
		proposal_last_response_at
	) VALUES (?, ?, 1, ?)
	ON DUPLICATE KEY UPDATE
		proposal_avg_response_time_seconds = CASE
			WHEN owner_response_metrics.proposal_avg_response_time_seconds IS NULL OR owner_response_metrics.proposal_total_responses = 0 THEN VALUES(proposal_avg_response_time_seconds)
			ELSE FLOOR((owner_response_metrics.proposal_avg_response_time_seconds * owner_response_metrics.proposal_total_responses + VALUES(proposal_avg_response_time_seconds)) / (owner_response_metrics.proposal_total_responses + 1))
		END,
		proposal_total_responses = owner_response_metrics.proposal_total_responses + 1,
		proposal_last_response_at = GREATEST(COALESCE(owner_response_metrics.proposal_last_response_at, VALUES(proposal_last_response_at)), VALUES(proposal_last_response_at))`

	defer a.ObserveOnComplete("upsert_proposal_response", query)()

	if _, execErr := a.ExecContext(ctx, tx, "upsert_proposal_response", query,
		input.OwnerID,
		input.DeltaSeconds,
		input.RespondedAt,
	); execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.owner_metrics.proposal.upsert_error", "owner_id", input.OwnerID, "err", execErr)
		return fmt.Errorf("upsert owner proposal metrics: %w", execErr)
	}

	return nil
}
