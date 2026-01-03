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
// Purpose:
//   - Provide a single source of truth for DB → Domain mapping for all session queries
//   - Ensure consistent NULL handling across all read paths
//   - Reduce duplication in per-method files while keeping methods concise (project guide)
//
// Flow:
//  1. Scan row into SessionEntity (with sql.Null* fields to preserve NULL semantics)
//  2. Convert entity to domain model using SessionEntityToDomain
//
// Parameters:
//   - ctx: Context carrying tracing/logging
//   - scanner: Implements Scan (either *sql.Row or *sql.Rows)
//   - operation: Operation name for logging/metrics (e.g., "get_by_id")
//
// Returns:
//   - session: SessionInterface populated from row data
//   - error: sql.ErrNoRows when no data; wrapped infrastructure errors otherwise
//
// Error Handling:
//   - sql.ErrNoRows is returned verbatim for service-layer mapping (404/empty)
//   - Infrastructure scan errors are logged, traced, and wrapped with operation context
//
// Usage:
//   - Used by all read methods (get by id, get by hash, list by user) to keep behavior uniform
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
