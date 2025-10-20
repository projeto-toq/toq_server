package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/entity"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *ScheduleAdapter) GetAgendaByListingID(ctx context.Context, tx *sql.Tx, listingID int64) (schedulemodel.AgendaInterface, error) {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return nil, err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, listing_id, owner_id, timezone FROM listing_agendas WHERE listing_id = ? LIMIT 1`

	row := exec.QueryRowContext(ctx, query, listingID)

	var agendaEntity entity.AgendaEntity
	if err = row.Scan(&agendaEntity.ID, &agendaEntity.ListingID, &agendaEntity.OwnerID, &agendaEntity.Timezone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.get_agenda.scan_error", "listing_id", listingID, "err", err)
		return nil, fmt.Errorf("scan agenda by listing id: %w", err)
	}

	return converters.ToAgendaModel(agendaEntity), nil
}
