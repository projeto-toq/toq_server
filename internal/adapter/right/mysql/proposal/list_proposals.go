package mysqlproposaladapter

import (
	"context"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/entities"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListProposals returns paginated proposals with summary counters.
func (a *ProposalAdapter) ListProposals(ctx context.Context, filter proposalmodel.ListFilter) (proposalmodel.ListResult, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return proposalmodel.ListResult{}, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	whereClause, args := buildProposalFilters(filter)
	limit, offset := normalizePagination(filter.Page, filter.Limit)
	orderBy := normalizeOrderBy(filter.SortBy)
	orderDirection := "DESC"
	if strings.EqualFold(filter.SortOrder, "asc") {
		orderDirection = "ASC"
	}

	query := fmt.Sprintf(`SELECT
		id,
		listing_identity_id,
		realtor_id,
		owner_id,
		transaction_type,
		payment_method,
		proposed_value,
		original_value,
		down_payment,
		installments,
		accepts_exchange,
		rental_months,
		guarantee_type,
		security_deposit,
		client_name,
		client_phone,
		proposal_notes,
		owner_notes,
		rejection_reason,
		status,
		expires_at,
		accepted_at,
		rejected_at,
		cancelled_at,
		is_favorite,
		created_at,
		updated_at,
		deleted
	FROM proposals
	%s
	ORDER BY %s %s
	LIMIT ? OFFSET ?`, whereClause, orderBy, orderDirection)

	rows, queryErr := a.QueryContext(ctx, nil, "list_proposals", query, append(args, limit, offset)...)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.proposal.list.query_error", "err", queryErr)
		return proposalmodel.ListResult{}, fmt.Errorf("list proposals: %w", queryErr)
	}
	defer rows.Close()

	items := make([]proposalmodel.ProposalInterface, 0)
	for rows.Next() {
		entity, scanErr := scanProposalEntity(rows)
		if scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.proposal.list.scan_error", "err", scanErr)
			return proposalmodel.ListResult{}, fmt.Errorf("scan proposal: %w", scanErr)
		}
		items = append(items, converters.ToProposalModel(entity))
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.proposal.list.rows_error", "err", rowsErr)
		return proposalmodel.ListResult{}, fmt.Errorf("iterate proposals: %w", rowsErr)
	}

	total, countErr := a.countProposals(ctx, whereClause, args)
	if countErr != nil {
		return proposalmodel.ListResult{}, countErr
	}

	summary, summaryErr := a.aggregateProposalSummary(ctx, whereClause, args)
	if summaryErr != nil {
		return proposalmodel.ListResult{}, summaryErr
	}

	return proposalmodel.ListResult{
		Items:   items,
		Total:   total,
		Summary: summary,
	}, nil
}

// countProposals returns the total rows matching filters.
func (a *ProposalAdapter) countProposals(ctx context.Context, whereClause string, args []interface{}) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM proposals %s", whereClause)
	var total int64
	row := a.QueryRowContext(ctx, nil, "count_proposals", query, args...)
	if err := row.Scan(&total); err != nil {
		utils.SetSpanError(ctx, err)
		logger := utils.LoggerFromContext(ctx)
		logger.Error("mysql.proposal.count.scan_error", "err", err)
		return 0, fmt.Errorf("count proposals: %w", err)
	}
	return total, nil
}

// aggregateProposalSummary calculates status counters and monetary extremes for the filtered set.
func (a *ProposalAdapter) aggregateProposalSummary(ctx context.Context, whereClause string, args []interface{}) (proposalmodel.Stats, error) {
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
	%s`, whereClause)

	stats := proposalmodel.Stats{}
	row := a.QueryRowContext(ctx, nil, "aggregate_proposal_stats", query, args...)
	if err := row.Scan(
		&stats.TotalProposals,
		&stats.PendingCount,
		&stats.AcceptedCount,
		&stats.RejectedCount,
		&stats.CancelledCount,
		&stats.ExpiredCount,
		&stats.HighestProposal,
		&stats.LowestProposal,
		&stats.AverageProposal,
	); err != nil {
		utils.SetSpanError(ctx, err)
		logger := utils.LoggerFromContext(ctx)
		logger.Error("mysql.proposal.stats.scan_error", "err", err)
		return proposalmodel.Stats{}, fmt.Errorf("aggregate proposal stats: %w", err)
	}

	return stats, nil
}

// scanProposalEntity maps a row into a ProposalEntity.
func scanProposalEntity(scanner interface{ Scan(dest ...any) error }) (entities.ProposalEntity, error) {
	entity := entities.ProposalEntity{}
	if err := scanner.Scan(
		&entity.ID,
		&entity.ListingIdentityID,
		&entity.RealtorID,
		&entity.OwnerID,
		&entity.TransactionType,
		&entity.PaymentMethod,
		&entity.ProposedValue,
		&entity.OriginalValue,
		&entity.DownPayment,
		&entity.Installments,
		&entity.AcceptsExchange,
		&entity.RentalMonths,
		&entity.GuaranteeType,
		&entity.SecurityDeposit,
		&entity.ClientName,
		&entity.ClientPhone,
		&entity.ProposalNotes,
		&entity.OwnerNotes,
		&entity.RejectionReason,
		&entity.Status,
		&entity.ExpiresAt,
		&entity.AcceptedAt,
		&entity.RejectedAt,
		&entity.CancelledAt,
		&entity.IsFavorite,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.Deleted,
	); err != nil {
		return entities.ProposalEntity{}, err
	}
	return entity, nil
}

func normalizePagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	return limit, offset
}

func normalizeOrderBy(sortBy string) string {
	switch strings.ToLower(sortBy) {
	case "proposedvalue":
		return "proposed_value"
	case "expiresat":
		return "expires_at"
	case "status":
		return "status"
	default:
		return "created_at"
	}
}

func buildProposalFilters(filter proposalmodel.ListFilter) (string, []interface{}) {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)

	if !filter.IncludeDeleted {
		conditions = append(conditions, "deleted = FALSE")
	}

	if len(filter.Statuses) > 0 {
		placeholders := buildPlaceholders(len(filter.Statuses))
		conditions = append(conditions, fmt.Sprintf("status IN (%s)", placeholders))
		for _, status := range filter.Statuses {
			args = append(args, string(status))
		}
	}

	if filter.ListingIdentityID != nil {
		conditions = append(conditions, "listing_identity_id = ?")
		args = append(args, *filter.ListingIdentityID)
	}

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

	if filter.MinValue != nil {
		conditions = append(conditions, "proposed_value >= ?")
		args = append(args, *filter.MinValue)
	}

	if filter.MaxValue != nil {
		conditions = append(conditions, "proposed_value <= ?")
		args = append(args, *filter.MaxValue)
	}

	if len(filter.TransactionTypes) > 0 {
		placeholders := buildPlaceholders(len(filter.TransactionTypes))
		conditions = append(conditions, fmt.Sprintf("transaction_type IN (%s)", placeholders))
		for _, t := range filter.TransactionTypes {
			args = append(args, string(t))
		}
	}

	if len(filter.PaymentMethods) > 0 {
		placeholders := buildPlaceholders(len(filter.PaymentMethods))
		conditions = append(conditions, fmt.Sprintf("payment_method IN (%s)", placeholders))
		for _, p := range filter.PaymentMethods {
			args = append(args, string(p))
		}
	}

	if len(conditions) == 0 {
		return "", args
	}

	return "WHERE " + strings.Join(conditions, " AND "), args
}

func buildPlaceholders(n int) string {
	if n <= 0 {
		return ""
	}

	placeholders := make([]string, n)
	for i := 0; i < n; i++ {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ",")
}
