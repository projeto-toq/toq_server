package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/entity"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PhotoSessionAdapter) ListAvailableSlots(ctx context.Context, tx *sql.Tx, params photosessionmodel.SlotListParams) (slots []photosessionmodel.PhotographerSlotInterface, total int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	filters := []string{"status = 'AVAILABLE'"}
	args := make([]any, 0)

	if params.From != nil {
		filters = append(filters, "slot_date >= ?")
		args = append(args, params.From)
	}

	if params.To != nil {
		filters = append(filters, "slot_date <= ?")
		args = append(args, params.To)
	}

	if params.Period != nil {
		filters = append(filters, "period = ?")
		args = append(args, string(*params.Period))
	}

	whereClause := strings.Join(filters, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM photographer_time_slots WHERE %s", whereClause)

	countStmt, err := tx.PrepareContext(ctx, countQuery)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.list_slots.prepare_count_error", "err", err)
		return nil, 0, fmt.Errorf("prepare count photographer slots: %w", err)
	}
	defer countStmt.Close()

	countArgs := append([]any{}, args...)

	if err = countStmt.QueryRowContext(ctx, countArgs...).Scan(&total); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.list_slots.count_error", "err", err)
		return nil, 0, fmt.Errorf("count photographer slots: %w", err)
	}

	sortColumn := params.SortColumn
	if sortColumn == "" {
		sortColumn = "slot_date"
	}
	sortDirection := params.SortDirection
	if sortDirection == "" {
		sortDirection = "ASC"
	}
	orderClause := fmt.Sprintf("%s %s", sortColumn, sortDirection)
	if sortColumn != "slot_date" {
		orderClause += ", slot_date ASC"
	}
	orderClause += ", period ASC"

	query := fmt.Sprintf(`
		SELECT id, photographer_user_id, slot_date, period, status, reservation_token, reserved_until, booked_at
		FROM photographer_time_slots
		WHERE %s
		ORDER BY %s
		LIMIT ? OFFSET ?
	`, whereClause, orderClause)

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.list_slots.prepare_query_error", "err", err)
		return nil, 0, fmt.Errorf("prepare list photographer slots: %w", err)
	}
	defer stmt.Close()

	queryArgs := append(args, params.Limit, params.Offset)

	rows, err := stmt.QueryContext(ctx, queryArgs...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.list_slots.query_error", "err", err)
		return nil, 0, fmt.Errorf("query photographer slots: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			reservationToken sql.NullString
			reservedUntil    sql.NullTime
			bookedAt         sql.NullTime
			entitySlot       entity.SlotEntity
		)

		if err = rows.Scan(
			&entitySlot.ID,
			&entitySlot.PhotographerUserID,
			&entitySlot.SlotDate,
			&entitySlot.Period,
			&entitySlot.Status,
			&reservationToken,
			&reservedUntil,
			&bookedAt,
		); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.photo_session.list_slots.scan_error", "err", err)
			return nil, 0, fmt.Errorf("scan photographer slot: %w", err)
		}

		if reservationToken.Valid {
			token := reservationToken.String
			entitySlot.ReservationToken = &token
		}
		if reservedUntil.Valid {
			value := reservedUntil.Time
			entitySlot.ReservedUntil = &value
		}
		if bookedAt.Valid {
			value := bookedAt.Time
			entitySlot.BookedAt = &value
		}

		slots = append(slots, converters.ToSlotModel(entitySlot))
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.list_slots.rows_error", "err", err)
		return nil, 0, fmt.Errorf("iterate photographer slots: %w", err)
	}

	return slots, total, nil
}
