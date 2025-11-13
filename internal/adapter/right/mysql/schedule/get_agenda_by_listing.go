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

func (a *ScheduleAdapter) GetAgendaByListingIdentityID(ctx context.Context, tx *sql.Tx, listingIdentityID int64) (schedulemodel.AgendaInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, listing_identity_id, owner_id, timezone FROM listing_agendas WHERE listing_identity_id = ? LIMIT 1`

	row := a.QueryRowContext(ctx, tx, "select", query, listingIdentityID)

	var agendaEntity entity.AgendaEntity
	if err = row.Scan(&agendaEntity.ID, &agendaEntity.ListingIdentityID, &agendaEntity.OwnerID, &agendaEntity.Timezone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.get_agenda.scan_error", "listing_identity_id", listingIdentityID, "err", err)
		return nil, fmt.Errorf("scan agenda by listing id: %w", err)
	}

	return converters.ToAgendaModel(agendaEntity), nil
}
