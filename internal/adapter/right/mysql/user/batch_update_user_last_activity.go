package mysqluseradapter

import (
	"context"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) BatchUpdateUserLastActivity(ctx context.Context, userIDs []int64, timestamps []int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if len(userIDs) == 0 {
		return nil
	}

	if len(userIDs) != len(timestamps) {
		errLength := fmt.Errorf("userIDs and timestamps length mismatch")
		utils.SetSpanError(ctx, errLength)
		logger.Error("mysql.user.batch_update_last_activity.length_mismatch", "user_ids", len(userIDs), "timestamps", len(timestamps), "error", errLength)
		return errLength
	}

	// Build batch update query using CASE WHEN for better performance
	var queryBuilder strings.Builder
	queryBuilder.WriteString("UPDATE users SET last_activity_at = CASE id ")

	args := make([]interface{}, 0, len(userIDs)*2)

	for i, userID := range userIDs {
		queryBuilder.WriteString("WHEN ? THEN FROM_UNIXTIME(?) ")
		args = append(args, userID, timestamps[i])
	}

	queryBuilder.WriteString("ELSE last_activity_at END WHERE id IN (")

	// Add placeholders for WHERE IN clause
	placeholders := make([]string, len(userIDs))
	for i, userID := range userIDs {
		placeholders[i] = "?"
		args = append(args, userID)
	}

	queryBuilder.WriteString(strings.Join(placeholders, ","))
	queryBuilder.WriteString(")")

	query := queryBuilder.String()

	result, execErr := ua.ExecContext(ctx, nil, "update", query, args...)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.batch_update_last_activity.exec_error", "error", execErr)
		return fmt.Errorf("exec batch update: %w", execErr)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		logger.Warn("mysql.user.batch_update_last_activity.rows_affected_warning", "error", err)
	} else {
		logger.Debug("mysql.user.batch_update_last_activity.success", "affected_rows", affected, "batch_size", len(userIDs))
	}

	return
}
