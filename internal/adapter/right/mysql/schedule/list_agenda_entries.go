package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/entity"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const agendaDetailMaxPageSize = 100

func (a *ScheduleAdapter) ListAgendaEntries(ctx context.Context, tx *sql.Tx, filter schedulemodel.AgendaDetailFilter) (schedulemodel.AgendaEntriesPage, error) {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return schedulemodel.AgendaEntriesPage{}, err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	conditions := []string{"a.owner_id = ?", "a.listing_id = ?"}
	args := []any{filter.OwnerID, filter.ListingID}

	if !filter.Range.From.IsZero() {
		conditions = append(conditions, "e.ends_at > ?")
		args = append(args, filter.Range.From)
	}
	if !filter.Range.To.IsZero() {
		conditions = append(conditions, "e.starts_at < ?")
		args = append(args, filter.Range.To)
	}

	where := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM listing_agenda_entries e INNER JOIN listing_agendas a ON a.id = e.agenda_id WHERE %s", where)

	var total int64
	if err = exec.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.agenda_detail.count_error", "listing_id", filter.ListingID, "err", err)
		return schedulemodel.AgendaEntriesPage{}, fmt.Errorf("count agenda entries: %w", err)
	}

	limit, offset := defaultPagination(filter.Pagination.Limit, filter.Pagination.Page, agendaDetailMaxPageSize)

	query := fmt.Sprintf(`
		SELECT e.id, e.agenda_id, e.entry_type, e.starts_at, e.ends_at, e.blocking, e.reason, e.visit_id, e.photo_booking_id
		FROM listing_agenda_entries e
		INNER JOIN listing_agendas a ON a.id = e.agenda_id
		WHERE %s
		ORDER BY e.starts_at
		LIMIT ? OFFSET ?
	`, where)

	argsWithPagination := append(append([]any{}, args...), limit, offset)

	rows, err := exec.QueryContext(ctx, query, argsWithPagination...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.agenda_detail.query_error", "listing_id", filter.ListingID, "err", err)
		return schedulemodel.AgendaEntriesPage{}, fmt.Errorf("query agenda entries: %w", err)
	}
	defer rows.Close()

	entries := make([]schedulemodel.AgendaEntryInterface, 0)
	for rows.Next() {
		var entryEntity entity.EntryEntity
		if err = rows.Scan(&entryEntity.ID, &entryEntity.AgendaID, &entryEntity.EntryType, &entryEntity.StartsAt, &entryEntity.EndsAt, &entryEntity.Blocking, &entryEntity.Reason, &entryEntity.VisitID, &entryEntity.PhotoBookingID); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.schedule.agenda_detail.scan_error", "listing_id", filter.ListingID, "err", err)
			return schedulemodel.AgendaEntriesPage{}, fmt.Errorf("scan agenda entry: %w", err)
		}
		entries = append(entries, converters.ToEntryModel(entryEntity))
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.agenda_detail.rows_error", "listing_id", filter.ListingID, "err", err)
		return schedulemodel.AgendaEntriesPage{}, fmt.Errorf("iterate agenda entries: %w", err)
	}

	return schedulemodel.AgendaEntriesPage{Entries: entries, Total: total}, nil
}
