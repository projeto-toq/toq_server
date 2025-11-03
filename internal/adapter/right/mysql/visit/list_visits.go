package mysqlvisitadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/visit/converters"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const visitsMaxPageSize = 50

func (a *VisitAdapter) ListVisits(ctx context.Context, tx *sql.Tx, filter listingmodel.VisitListFilter) (listingmodel.VisitListResult, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return listingmodel.VisitListResult{}, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	conditions := make([]string, 0)
	args := make([]any, 0)

	if filter.ListingID != nil {
		conditions = append(conditions, "listing_id = ?")
		args = append(args, *filter.ListingID)
	}
	if filter.OwnerID != nil {
		conditions = append(conditions, "owner_id = ?")
		args = append(args, *filter.OwnerID)
	}
	if filter.RealtorID != nil {
		conditions = append(conditions, "realtor_id = ?")
		args = append(args, *filter.RealtorID)
	}
	if len(filter.Statuses) > 0 {
		placeholders := make([]string, len(filter.Statuses))
		for i, status := range filter.Statuses {
			placeholders[i] = "?"
			args = append(args, string(status))
		}
		conditions = append(conditions, fmt.Sprintf("status IN (%s)", strings.Join(placeholders, ",")))
	}
	if filter.From != nil {
		conditions = append(conditions, "scheduled_end >= ?")
		args = append(args, *filter.From)
	}
	if filter.To != nil {
		conditions = append(conditions, "scheduled_start <= ?")
		args = append(args, *filter.To)
	}

	where := "1=1"
	if len(conditions) > 0 {
		where = strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM listing_visits WHERE %s", where)
	var total int64
	countRow := a.QueryRowContext(ctx, tx, "list_visits_count", countQuery, args...)
	if err = countRow.Scan(&total); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.visit.list.count_error", "err", err)
		return listingmodel.VisitListResult{}, fmt.Errorf("count visits: %w", err)
	}

	limit, offset := defaultPagination(filter.Limit, filter.Page, visitsMaxPageSize)

	query := fmt.Sprintf(`
		SELECT id, listing_id, owner_id, realtor_id, scheduled_start, scheduled_end, status, cancel_reason, notes, created_by, updated_by
		FROM listing_visits
		WHERE %s
		ORDER BY scheduled_start ASC
		LIMIT ? OFFSET ?
	`, where)

	params := append(make([]any, 0, len(args)+2), args...)
	params = append(params, limit, offset)

	rows, err := a.QueryContext(ctx, tx, "list_visits", query, params...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.visit.list.query_error", "err", err)
		return listingmodel.VisitListResult{}, fmt.Errorf("query visits: %w", err)
	}
	defer rows.Close()

	visits := make([]listingmodel.VisitInterface, 0)
	for rows.Next() {
		visitEntity, scanErr := scanVisitEntity(rows)
		if scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.visit.list.scan_error", "err", scanErr)
			return listingmodel.VisitListResult{}, fmt.Errorf("scan visit: %w", scanErr)
		}
		visits = append(visits, converters.ToVisitModel(visitEntity))
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.visit.list.rows_error", "err", err)
		return listingmodel.VisitListResult{}, fmt.Errorf("iterate visits: %w", err)
	}

	return listingmodel.VisitListResult{Visits: visits, Total: total}, nil
}
