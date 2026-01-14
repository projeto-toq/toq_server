package mysqlproposaladapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/entities"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListOwnerSummaries returns aggregated owner information for the provided IDs, including response metrics.
func (a *ProposalAdapter) ListOwnerSummaries(
	ctx context.Context,
	tx *sql.Tx,
	ownerIDs []int64,
) ([]proposalmodel.OwnerSummary, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if len(ownerIDs) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(ownerIDs))
	args := make([]interface{}, len(ownerIDs))
	for i, id := range ownerIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`SELECT
        u.id AS owner_id,
        u.full_name,
        TIMESTAMPDIFF(MONTH, COALESCE(u.created_at, u.last_activity_at), UTC_TIMESTAMP()) AS member_since_months,
        orm.proposal_avg_response_time_seconds,
        orm.visit_avg_response_time_seconds
    FROM users u
    LEFT JOIN owner_response_metrics orm ON orm.user_id = u.id
    WHERE u.id IN (%s)`, strings.Join(placeholders, ","))

	rows, queryErr := a.QueryContext(ctx, tx, "list_owner_summaries", query, args...)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.proposal.owner_summary.query_error", "err", queryErr)
		return nil, fmt.Errorf("list owner summaries: %w", queryErr)
	}
	defer rows.Close()

	summaries := make([]proposalmodel.OwnerSummary, 0, len(ownerIDs))
	for rows.Next() {
		entity := entities.OwnerSummaryEntity{}
		if scanErr := rows.Scan(
			&entity.OwnerID,
			&entity.FullName,
			&entity.MemberSinceMonths,
			&entity.ProposalAvgSeconds,
			&entity.VisitAvgSeconds,
		); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.proposal.owner_summary.scan_error", "err", scanErr)
			return nil, fmt.Errorf("scan owner summary: %w", scanErr)
		}
		summaries = append(summaries, converters.ToOwnerSummaryModel(entity))
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.proposal.owner_summary.rows_error", "err", rowsErr)
		return nil, fmt.Errorf("iterate owner summaries: %w", rowsErr)
	}

	return summaries, nil
}
