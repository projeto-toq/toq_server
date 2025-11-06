package mysqluseradapter

import (
	"context"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// BatchUpdateUserLastActivity updates last_activity_at for multiple users in a single query
//
// This function performs a batch update using MySQL's CASE WHEN pattern, which is significantly
// more efficient than executing N individual UPDATE statements. All updates are performed
// within a single transaction.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - userIDs: Slice of user IDs to update
//   - timestamps: Unix timestamps (seconds) corresponding to each user ID
//
// Returns:
//   - error: Length mismatch errors, database errors
//
// Performance:
//   - Single query updates all users (massive performance gain vs N queries)
//   - No transaction required (atomic at statement level)
//   - CASE WHEN pattern ensures correct timestamp per user
//
// Query Pattern:
//
//	UPDATE users SET last_activity_at = CASE id
//	  WHEN 123 THEN FROM_UNIXTIME(1699300000)
//	  WHEN 456 THEN FROM_UNIXTIME(1699300001)
//	  ELSE last_activity_at
//	END WHERE id IN (123, 456)
//
// Business Rules:
//   - userIDs and timestamps must have same length
//   - Empty slices are allowed (early return, no-op)
//   - WHERE IN clause ensures only specified users are affected
func (ua *UserAdapter) BatchUpdateUserLastActivity(ctx context.Context, userIDs []int64, timestamps []int64) (err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Early return for empty batch (optimization - avoid unnecessary query)
	if len(userIDs) == 0 {
		return nil
	}

	// Validate input arrays have matching lengths (defensive programming)
	if len(userIDs) != len(timestamps) {
		errLength := fmt.Errorf("userIDs and timestamps length mismatch")
		utils.SetSpanError(ctx, errLength)
		logger.Error("mysql.user.batch_update_last_activity.length_mismatch", "user_ids", len(userIDs), "timestamps", len(timestamps), "error", errLength)
		return errLength
	}

	// Build batch update using CASE WHEN pattern for single-query efficiency
	// This approach updates multiple rows with different values in one transaction
	// instead of N individual UPDATE statements (massive performance gain)
	//
	// Generated SQL example for 3 users:
	// UPDATE users SET last_activity_at = CASE id
	//   WHEN 123 THEN FROM_UNIXTIME(1699300000)
	//   WHEN 456 THEN FROM_UNIXTIME(1699300001)
	//   WHEN 789 THEN FROM_UNIXTIME(1699300002)
	//   ELSE last_activity_at
	// END WHERE id IN (123, 456, 789)
	var queryBuilder strings.Builder
	queryBuilder.WriteString("UPDATE users SET last_activity_at = CASE id ")

	args := make([]interface{}, 0, len(userIDs)*2)

	// Build CASE WHEN clauses (one per user)
	for i, userID := range userIDs {
		queryBuilder.WriteString("WHEN ? THEN FROM_UNIXTIME(?) ")
		args = append(args, userID, timestamps[i])
	}

	// ELSE clause preserves unchanged rows (safety net)
	queryBuilder.WriteString("ELSE last_activity_at END WHERE id IN (")

	// Add WHERE IN clause for safety (only update specified users)
	placeholders := make([]string, len(userIDs))
	for i, userID := range userIDs {
		placeholders[i] = "?"
		args = append(args, userID)
	}

	queryBuilder.WriteString(strings.Join(placeholders, ","))
	queryBuilder.WriteString(")")

	query := queryBuilder.String()

	// Execute batch update (no transaction - read-modify-write not needed)
	result, execErr := ua.ExecContext(ctx, nil, "update", query, args...)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.batch_update_last_activity.exec_error", "error", execErr)
		return fmt.Errorf("exec batch update: %w", execErr)
	}

	// Log success metrics for observability (debug level for high-frequency operations)
	affected, err := result.RowsAffected()
	if err != nil {
		logger.Warn("mysql.user.batch_update_last_activity.rows_affected_warning", "error", err)
	} else {
		logger.Debug("mysql.user.batch_update_last_activity.success", "affected_rows", affected, "batch_size", len(userIDs))
	}

	return
}
