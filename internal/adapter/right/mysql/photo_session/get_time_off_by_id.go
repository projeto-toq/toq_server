package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/entity"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetTimeOffByID retrieves a photographer time-off entry by identifier.
func (a *PhotoSessionAdapter) GetTimeOffByID(ctx context.Context, tx *sql.Tx, timeOffID uint64) (photosessionmodel.PhotographerTimeOffInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
        SELECT id, photographer_user_id, start_date, end_date, reason
        FROM photographer_time_off
        WHERE id = ?
    `

	row := tx.QueryRowContext(ctx, query, timeOffID)

	var ent entity.TimeOffEntity
	if err = row.Scan(&ent.ID, &ent.PhotographerUserID, &ent.StartDate, &ent.EndDate, &ent.Reason); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.get_time_off.scan_error", "time_off_id", timeOffID, "err", err)
		return nil, fmt.Errorf("get photographer time off: %w", err)
	}

	return converters.ToTimeOffModel(ent), nil
}
