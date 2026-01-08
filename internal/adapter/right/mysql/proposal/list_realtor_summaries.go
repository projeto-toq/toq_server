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

// ListRealtorSummaries returns aggregated realtor information for the provided IDs.
func (a *ProposalAdapter) ListRealtorSummaries(
	ctx context.Context,
	tx *sql.Tx,
	realtorIDs []int64,
) ([]proposalmodel.RealtorSummary, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if len(realtorIDs) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(realtorIDs))
	args := make([]interface{}, len(realtorIDs))
	for i, id := range realtorIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`SELECT
		u.id AS realtor_id,
		u.full_name,
		u.nick_name,
		TIMESTAMPDIFF(MONTH, COALESCE(u.created_at, u.last_activity_at), UTC_TIMESTAMP()) AS usage_months,
		COALESCE(stats.total_proposals, 0) AS proposals_count
	FROM users u
	LEFT JOIN (
		SELECT realtor_id, COUNT(*) AS total_proposals
		FROM proposals
		WHERE deleted = 0
		GROUP BY realtor_id
	) stats ON stats.realtor_id = u.id
	WHERE u.id IN (%s)`, strings.Join(placeholders, ","))

	rows, queryErr := a.QueryContext(ctx, tx, "list_realtor_summaries", query, args...)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.proposal.realtor_summary.query_error", "err", queryErr)
		return nil, fmt.Errorf("list realtor summaries: %w", queryErr)
	}
	defer rows.Close()

	summaries := make([]proposalmodel.RealtorSummary, 0, len(realtorIDs))

	for rows.Next() {
		entity := entities.RealtorSummaryEntity{}
		if scanErr := rows.Scan(
			&entity.RealtorID,
			&entity.FullName,
			&entity.NickName,
			&entity.UsageMonths,
			&entity.ProposalsCount,
		); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.proposal.realtor_summary.scan_error", "err", scanErr)
			return nil, fmt.Errorf("scan realtor summary: %w", scanErr)
		}
		summaries = append(summaries, converters.ToRealtorSummaryModel(entity))
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.proposal.realtor_summary.rows_error", "err", rowsErr)
		return nil, fmt.Errorf("iterate realtor summaries: %w", rowsErr)
	}

	return summaries, nil
}
