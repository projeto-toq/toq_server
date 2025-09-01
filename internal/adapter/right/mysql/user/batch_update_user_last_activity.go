package mysqluseradapter

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) BatchUpdateUserLastActivity(ctx context.Context, userIDs []int64, timestamps []int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	if len(userIDs) == 0 {
		return nil
	}

	if len(userIDs) != len(timestamps) {
		return fmt.Errorf("userIDs and timestamps length mismatch")
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

	// Execute batch update
	stmt, err := ua.db.DB.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("mysqluseradapter/BatchUpdateUserLastActivity: error preparing statement", "error", err)
		return fmt.Errorf("prepare batch update: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		slog.Error("mysqluseradapter/BatchUpdateUserLastActivity: error executing statement", "error", err)
		return fmt.Errorf("exec batch update: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		slog.Warn("mysqluseradapter/BatchUpdateUserLastActivity: could not get affected rows", "error", err)
	} else {
		slog.Debug("Batch updated user activities", "affected_rows", affected, "batch_size", len(userIDs))
	}

	return
}
