package sessionmysqladapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sessionconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/session/converters"
	sessionentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/session/entities"
	sessionmodel "github.com/projeto-toq/toq_server/internal/core/model/session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// rowScanner defines the interface for scanning database rows
// Satisfied by both *sql.Row and *sql.Rows
type rowScanner interface {
	Scan(dest ...any) error
}

// mapSessionFromScanner scans a database row into SessionEntity then converts to domain
//
// This function centralizes the row scanning logic and delegates conversion to the
// dedicated converter. This ensures a single source of truth for DB → Domain mapping.
//
// Flow:
//  1. Scan row into SessionEntity (with proper sql.Null* handling)
//  2. Convert entity to domain model via SessionEntityToDomain converter
//
// Parameters:
//   - ctx: Context for error logging and tracing
//   - scanner: Row or Rows interface implementing Scan()
//   - operation: Operation name for contextual error messages (e.g., "get_by_id", "list")
//
// Returns:
//   - session: SessionInterface populated from row data
//   - error: sql.ErrNoRows if no data, or scan/conversion errors
//
// Important:
//   - Returns sql.ErrNoRows directly (no wrapping) for repository-level handling
//   - Logs and marks span on scan errors (infrastructure failure)
//   - Uses SessionEntityToDomain for all NULL → value conversions
func (sa *SessionAdapter) mapSessionFromScanner(ctx context.Context, scanner rowScanner, operation string) (sessionmodel.SessionInterface, error) {
	// Prepare entity with sql.Null* types for proper NULL handling
	var entity sessionentities.SessionEntity
	var revokedInt int64 // MySQL TINYINT scanned as int64, converted to bool

	// Scan row into entity fields (order must match SELECT column order)
	if err := scanner.Scan(
		&entity.ID,
		&entity.UserID,
		&entity.RefreshHash,
		&entity.TokenJTI,
		&entity.ExpiresAt,
		&entity.AbsoluteExpiresAt,
		&entity.CreatedAt,
		&entity.RotatedAt,
		&entity.UserAgent,
		&entity.IP,
		&entity.DeviceID,
		&entity.RotationCounter,
		&entity.LastRefreshAt,
		&revokedInt,
	); err != nil {
		// Return sql.ErrNoRows directly (no wrapping) for service layer handling
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		// Infrastructure error: log and mark span
		utils.SetSpanError(ctx, err)
		logger := utils.LoggerFromContext(ctx)
		logger.Error(fmt.Sprintf("mysql.session.%s.scan_error", operation), "error", err)
		return nil, fmt.Errorf("scan session (%s): %w", operation, err)
	}

	// Convert TINYINT (0/1) to bool
	entity.Revoked = revokedInt == 1

	// Convert entity to domain using centralized converter
	// This ensures single source of truth for NULL handling logic
	session := sessionconverters.SessionEntityToDomain(entity)

	return session, nil
}
