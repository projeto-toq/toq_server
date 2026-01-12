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

// ListProposals returns paginated proposals scoped by actor filters.
func (a *ProposalAdapter) ListProposals(ctx context.Context, tx *sql.Tx, filter proposalmodel.ListFilter) (proposalmodel.ListResult, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return proposalmodel.ListResult{}, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	whereClause, args := buildListFilters(filter)
	limit, offset := normalizePagination(filter.Page, filter.Limit)

	query := fmt.Sprintf(`SELECT
		p.id,
		p.listing_identity_id,
		p.realtor_id,
		p.owner_id,
		p.proposal_text,
		p.rejection_reason,
		p.status,
		p.accepted_at,
		p.rejected_at,
		p.cancelled_at,
		p.first_owner_action_at,
		p.created_at,
		p.deleted,
		COALESCE(d.documents_count, 0) AS documents_count
	FROM proposals p
	LEFT JOIN (
		SELECT proposal_id, COUNT(*) AS documents_count
		FROM proposal_documents
		GROUP BY proposal_id
	) d ON d.proposal_id = p.id
	%s
	ORDER BY p.id DESC
	LIMIT ? OFFSET ?`, whereClause)

	rows, queryErr := a.QueryContext(ctx, tx, "list_proposals", query, append(args, limit, offset)...)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.proposal.list.query_error", "err", queryErr)
		return proposalmodel.ListResult{}, fmt.Errorf("list proposals: %w", queryErr)
	}
	defer rows.Close()

	items := make([]proposalmodel.ProposalInterface, 0)
	for rows.Next() {
		entity := entities.ProposalEntity{}
		if scanErr := rows.Scan(
			&entity.ID,
			&entity.ListingIdentityID,
			&entity.RealtorID,
			&entity.OwnerID,
			&entity.ProposalText,
			&entity.RejectionReason,
			&entity.Status,
			&entity.AcceptedAt,
			&entity.RejectedAt,
			&entity.CancelledAt,
			&entity.FirstOwnerAction,
			&entity.CreatedAt,
			&entity.Deleted,
			&entity.DocumentsCount,
		); scanErr != nil {
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

	total, countErr := a.countFilteredProposals(ctx, tx, whereClause, args)
	if countErr != nil {
		return proposalmodel.ListResult{}, countErr
	}

	return proposalmodel.ListResult{
		Items: items,
		Total: total,
	}, nil
}

func (a *ProposalAdapter) countFilteredProposals(ctx context.Context, tx *sql.Tx, whereClause string, args []interface{}) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM proposals p %s", whereClause)
	row := a.QueryRowContext(ctx, tx, "count_proposals", query, args...)
	var total int64
	if err := row.Scan(&total); err != nil {
		utils.SetSpanError(ctx, err)
		logger := utils.LoggerFromContext(ctx)
		logger.Error("mysql.proposal.count.scan_error", "err", err)
		return 0, fmt.Errorf("count proposals: %w", err)
	}
	return total, nil
}

func buildListFilters(filter proposalmodel.ListFilter) (string, []interface{}) {
	conditions := []string{"p.deleted = 0"}
	args := make([]interface{}, 0)

	switch filter.ActorScope {
	case proposalmodel.ActorScopeRealtor:
		conditions = append(conditions, "p.realtor_id = ?")
		args = append(args, filter.ActorID)
	case proposalmodel.ActorScopeOwner:
		conditions = append(conditions, "p.owner_id = ?")
		args = append(args, filter.ActorID)
	}

	if filter.ListingID != nil {
		conditions = append(conditions, "p.listing_identity_id = ?")
		args = append(args, *filter.ListingID)
	}

	if len(filter.Statuses) > 0 {
		placeholders := make([]string, len(filter.Statuses))
		for i, status := range filter.Statuses {
			placeholders[i] = "?"
			args = append(args, string(status))
		}
		conditions = append(conditions, fmt.Sprintf("p.status IN (%s)", strings.Join(placeholders, ",")))
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}
	return whereClause, args
}

func normalizePagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return limit, (page - 1) * limit
}
