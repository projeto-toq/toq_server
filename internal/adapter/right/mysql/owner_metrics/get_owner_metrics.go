package mysqlownermetricsadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/owner_metrics/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/owner_metrics/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetByOwnerID retrieves persisted SLA metrics for the given owner.
func (a *OwnerMetricsAdapter) GetByOwnerID(ctx context.Context, tx *sql.Tx, ownerID int64) (usermodel.OwnerResponseMetrics, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT
		user_id,
		visit_avg_response_time_seconds,
		visit_total_responses,
		visit_last_response_at,
		proposal_avg_response_time_seconds,
		proposal_total_responses,
		proposal_last_response_at
	FROM owner_response_metrics
	WHERE user_id = ?`

	row := a.QueryRowContext(ctx, tx, "get_owner_metrics", query, ownerID)
	entity := entities.OwnerMetricsEntity{}
	if scanErr := row.Scan(
		&entity.UserID,
		&entity.VisitAvgResponseSeconds,
		&entity.VisitTotalResponses,
		&entity.VisitLastResponseAt,
		&entity.ProposalAvgResponseSeconds,
		&entity.ProposalTotalResponses,
		&entity.ProposalLastResponseAt,
	); scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.owner_metrics.get.scan_error", "owner_id", ownerID, "err", scanErr)
		return nil, fmt.Errorf("get owner response metrics: %w", scanErr)
	}

	return converters.ToDomain(entity), nil
}
