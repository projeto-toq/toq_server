package sessionconverters

import (
	"database/sql"

	sessionentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/session/entities"
	sessionmodel "github.com/projeto-toq/toq_server/internal/core/model/session_model"
)

// SessionDomainToEntity converts a domain model to a database entity
//
// This converter handles the translation from clean domain types to database-specific
// types (sql.Null*), preparing data for database insertion/update.
//
// Conversion Rules:
//   - string → sql.NullString (Valid=true if non-empty, Valid=false if empty → NULL in DB)
//   - time.Time → sql.NullTime (Valid=true if not zero time, Valid=false if zero → NULL in DB)
//   - *time.Time → sql.NullTime (Valid=true if pointer non-nil, Valid=false if nil → NULL in DB)
//   - int → int (direct mapping; 0 is valid value)
//   - bool → bool (direct mapping; stored as TINYINT 0/1 in DB)
//
// NULL Semantics:
//   - TokenJTI: Empty string → NULL (session without access token)
//   - AbsoluteExpiresAt: Zero time → NULL (no absolute expiration)
//   - RotatedAt: nil pointer → NULL (never rotated)
//   - UserAgent/IP/DeviceID: Empty strings → NULL (metadata not available)
//   - RotationCounter: 0 is valid (initial state), NOT converted to NULL
//   - LastRefreshAt: nil pointer → NULL (never refreshed)
//
// Parameters:
//   - domain: SessionInterface from core layer
//
// Returns:
//   - entity: SessionEntity ready for database operations
//
// Important:
//   - ID may be 0 for new records (populated by AUTO_INCREMENT after INSERT)
//   - Empty strings are converted to NULL for optional VARCHAR/CHAR fields
//   - Zero times (IsZero()) are converted to NULL for optional DATETIME fields
//   - Nil pointers are converted to NULL for pointer time fields
//   - RotationCounter value 0 is preserved (not converted to NULL; matches DEFAULT 0)
//
// Example:
//
//	domain := sessionmodel.NewSession()
//	domain.SetUserID(456)
//	domain.SetRefreshHash("abc123...")
//	domain.SetExpiresAt(time.Now().Add(24 * time.Hour))
//	// TokenJTI left empty → will be NULL in DB
//	entity := SessionDomainToEntity(domain)
//	fmt.Println(entity.TokenJTI.Valid) // Output: false (NULL)
func SessionDomainToEntity(domain sessionmodel.SessionInterface) sessionentities.SessionEntity {
	entity := sessionentities.SessionEntity{}

	// Map mandatory fields (NOT NULL in schema)
	entity.ID = domain.GetID()
	entity.UserID = domain.GetUserID()
	entity.RefreshHash = domain.GetRefreshHash()
	entity.ExpiresAt = domain.GetExpiresAt()
	entity.CreatedAt = domain.GetCreatedAt()
	entity.RotationCounter = domain.GetRotationCounter() // 0 is valid, not NULL
	entity.Revoked = domain.GetRevoked()

	// Map optional fields - convert to sql.Null* with Valid based on value presence

	// TokenJTI: empty string → NULL (session without issued access token)
	tokenJTI := domain.GetTokenJTI()
	entity.TokenJTI = sql.NullString{
		String: tokenJTI,
		Valid:  tokenJTI != "",
	}

	// AbsoluteExpiresAt: zero time → NULL (no hard expiration limit)
	absoluteExp := domain.GetAbsoluteExpiresAt()
	entity.AbsoluteExpiresAt = sql.NullTime{
		Time:  absoluteExp,
		Valid: !absoluteExp.IsZero(),
	}

	// RotatedAt: nil pointer → NULL (session never refreshed)
	rotatedAt := domain.GetRotatedAt()
	if rotatedAt != nil {
		entity.RotatedAt = sql.NullTime{
			Time:  *rotatedAt,
			Valid: true,
		}
	} else {
		entity.RotatedAt = sql.NullTime{Valid: false}
	}

	// UserAgent: empty string → NULL (client didn't provide user-agent)
	userAgent := domain.GetUserAgent()
	entity.UserAgent = sql.NullString{
		String: userAgent,
		Valid:  userAgent != "",
	}

	// IP: empty string → NULL (IP not captured)
	ip := domain.GetIP()
	entity.IP = sql.NullString{
		String: ip,
		Valid:  ip != "",
	}

	// DeviceID: empty string → NULL (legacy session or device ID not provided)
	deviceID := domain.GetDeviceID()
	entity.DeviceID = sql.NullString{
		String: deviceID,
		Valid:  deviceID != "",
	}

	// LastRefreshAt: nil pointer → NULL (session never refreshed)
	lastRefresh := domain.GetLastRefreshAt()
	if lastRefresh != nil {
		entity.LastRefreshAt = sql.NullTime{
			Time:  *lastRefresh,
			Valid: true,
		}
	} else {
		entity.LastRefreshAt = sql.NullTime{Valid: false}
	}

	return entity
}
