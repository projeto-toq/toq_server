package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) DeleteAgencyRealtorRelation(ctx context.Context, tx *sql.Tx, agencyID int64, realtorID int64) (deleted int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM realtors_agency WHERE realtor_id = ? AND agency_id = ?;`

	result, execErr := ua.ExecContext(ctx, tx, "delete", query, realtorID, agencyID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.delete_agency_realtor_relation.exec_error", "error", execErr)
		return 0, fmt.Errorf("delete realtor-agency relation: %w", execErr)
	}

	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.user.delete_agency_realtor_relation.rows_affected_error", "error", rowsErr)
		return 0, fmt.Errorf("delete realtor-agency relation rows affected: %w", rowsErr)
	}

	if rowsAffected == 0 {
		return 0, sql.ErrNoRows
	}

	return rowsAffected, nil
}
