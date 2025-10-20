package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/schedule/converters"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *ScheduleAdapter) InsertAgenda(ctx context.Context, tx *sql.Tx, agenda schedulemodel.AgendaInterface) (uint64, error) {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return 0, err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity := converters.ToAgendaEntity(agenda)

	query := `INSERT INTO listing_agendas (listing_id, owner_id, timezone) VALUES (?, ?, ?)`
	result, err := exec.ExecContext(ctx, query, entity.ListingID, entity.OwnerID, entity.Timezone)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.insert_agenda.exec_error", "listing_id", entity.ListingID, "err", err)
		return 0, fmt.Errorf("insert agenda: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.insert_agenda.last_id_error", "listing_id", entity.ListingID, "err", err)
		return 0, fmt.Errorf("agenda last insert id: %w", err)
	}

	agenda.SetID(uint64(id))
	return uint64(id), nil
}
