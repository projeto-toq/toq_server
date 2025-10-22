package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	//"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/converters"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateTimeOff inserts a new time-off entry.
func (a *PhotoSessionAdapter) CreateTimeOff(ctx context.Context, tx *sql.Tx, timeOff photosessionmodel.PhotographerTimeOffInterface) (uint64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
		INSERT INTO photographer_time_off (photographer_user_id, start_date, end_date, reason)
		VALUES (?, ?, ?, ?)
	`

	res, execErr := tx.ExecContext(ctx, query,
		timeOff.PhotographerUserID(),
		timeOff.StartDate(),
		timeOff.EndDate(),
		timeOff.Reason(),
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.photo_session.create_time_off.exec_error", "err", execErr)
		return 0, fmt.Errorf("insert photographer time off: %w", execErr)
	}

	id, lastIDErr := res.LastInsertId()
	if lastIDErr != nil {
		utils.SetSpanError(ctx, lastIDErr)
		logger.Error("mysql.photo_session.create_time_off.last_insert_id_error", "err", lastIDErr)
		return 0, fmt.Errorf("photographer time off last insert id: %w", lastIDErr)
	}

	return uint64(id), nil
}
