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
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity := converters.ToAgendaEntity(agenda)

	query := `INSERT INTO listing_agendas (listing_identity_id, owner_id, timezone) VALUES (?, ?, ?)`
	result, execErr := a.ExecContext(ctx, tx, "insert", query, entity.ListingIdentityID, entity.OwnerID, entity.Timezone)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.schedule.insert_agenda.exec_error", "listing_identity_id", entity.ListingIdentityID, "err", execErr)
		return 0, fmt.Errorf("insert agenda: %w", execErr)
	}

	id, lastIDErr := result.LastInsertId()
	if lastIDErr != nil {
		utils.SetSpanError(ctx, lastIDErr)
		logger.Error("mysql.schedule.insert_agenda.last_id_error", "listing_identity_id", entity.ListingIdentityID, "err", lastIDErr)
		return 0, fmt.Errorf("agenda last insert id: %w", lastIDErr)
	}

	agenda.SetID(uint64(id))
	return uint64(id), nil
}
