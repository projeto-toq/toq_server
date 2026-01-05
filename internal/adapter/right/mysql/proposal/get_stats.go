package mysqlproposaladapter

import (
	"context"
	"fmt"
	"strings"

	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetStats aggregates proposal metrics for a scoped filter (realtor/owner/date).
func (a *ProposalAdapter) GetStats(ctx context.Context, filter proposalmodel.StatsFilter) (proposalmodel.Stats, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return proposalmodel.Stats{}, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	conditions := []string{"deleted = FALSE"}
	args := make([]interface{}, 0)

	if filter.RealtorID != nil {
		conditions = append(conditions, "realtor_id = ?")
		args = append(args, *filter.RealtorID)
	}

	if filter.OwnerID != nil {
		conditions = append(conditions, "owner_id = ?")
		args = append(args, *filter.OwnerID)
	}

	if filter.StartDate != nil {
		conditions = append(conditions, "created_at >= ?")
		args = append(args, *filter.StartDate)
	}

	if filter.EndDate != nil {
		conditions = append(conditions, "created_at <= ?")
		args = append(args, *filter.EndDate)
	}

	where := "WHERE " + strings.Join(conditions, " AND ")
	query := fmt.Sprintf(`SELECT
		COUNT(*) AS total_proposals,
		SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END) AS pending_count,
		SUM(CASE WHEN status = 'accepted' THEN 1 ELSE 0 END) AS accepted_count,
		SUM(CASE WHEN status = 'rejected' THEN 1 ELSE 0 END) AS rejected_count,
		SUM(CASE WHEN status = 'cancelled' THEN 1 ELSE 0 END) AS cancelled_count,
		SUM(CASE WHEN status = 'expired' THEN 1 ELSE 0 END) AS expired_count,
		MAX(proposed_value) AS highest_proposal,
		MIN(proposed_value) AS lowest_proposal,
		AVG(proposed_value) AS average_proposal
	FROM proposals
	%s`, where)

	stats := proposalmodel.Stats{}
	row := a.QueryRowContext(ctx, nil, "proposal_get_stats", query, args...)
	if scanErr := row.Scan(
		&stats.TotalProposals,
		&stats.PendingCount,
		&stats.AcceptedCount,
		&stats.RejectedCount,
		&stats.CancelledCount,
		&stats.ExpiredCount,
		&stats.HighestProposal,
		&stats.LowestProposal,
		&stats.AverageProposal,
	); scanErr != nil {
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.proposal.stats.scan_error", "err", scanErr)
		return proposalmodel.Stats{}, fmt.Errorf("get proposal stats: %w", scanErr)
	}

	return stats, nil
}
