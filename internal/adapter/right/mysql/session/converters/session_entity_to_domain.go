package sessionconverters

import (
	sessionentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/session/entities"
	sessionmodel "github.com/projeto-toq/toq_server/internal/core/model/session_model"
)

// SessionEntityToDomain converts a database entity to a domain model
//
// This converter handles the translation from database-specific types (sql.Null*)
// to clean domain types, ensuring the core layer remains decoupled from database concerns.
//
// Conversion Rules:
//   - sql.NullString → string (empty string if NULL or !Valid)
//   - sql.NullTime → time.Time (zero time if NULL; pointer types remain nil if NULL)
//   - Direct types → Direct types (id, user_id, refresh_hash, expires_at, created_at, rotation_counter, revoked)
//   - bool → bool (direct mapping; DB stores 0/1 as TINYINT)
//
// NULL Semantics:
//   - TokenJTI: Empty string for NULL (session without issued access token)
//   - AbsoluteExpiresAt: Zero time for NULL (no absolute expiration)
//   - RotatedAt: Pointer remains nil if NULL (never rotated)
//   - UserAgent/IP/DeviceID: Empty strings for NULL (metadata not captured)
//   - LastRefreshAt: Pointer remains nil if NULL (never refreshed)
//
// Parameters:
//   - entity: SessionEntity from database query
//
// Returns:
//   - session: SessionInterface with all fields populated from entity
//
// Example:
//
//	entity := sessionentities.SessionEntity{
//	    ID: 123,
//	    UserID: 456,
//	    RefreshHash: "abc123...",
//	    TokenJTI: sql.NullString{String: "uuid-here", Valid: true},
//	    // ... other fields
//	}
//	domain := SessionEntityToDomain(entity)
//	fmt.Println(domain.GetID()) // Output: 123
func SessionEntityToDomain(entity sessionentities.SessionEntity) sessionmodel.SessionInterface {
	session := sessionmodel.NewSession()

	// Map mandatory fields (NOT NULL in schema)
	session.SetID(entity.ID)
	session.SetUserID(entity.UserID)
	session.SetRefreshHash(entity.RefreshHash)
	session.SetExpiresAt(entity.ExpiresAt)
	session.SetCreatedAt(entity.CreatedAt)
	session.SetRotationCounter(entity.RotationCounter)
	session.SetRevoked(entity.Revoked)

	// Map optional fields (NULL in schema) - check Valid before accessing
	// TokenJTI: empty string if NULL (session created but no access token issued yet)
	if entity.TokenJTI.Valid {
		session.SetTokenJTI(entity.TokenJTI.String)
	}

	// AbsoluteExpiresAt: zero time if NULL (no hard expiration limit)
	if entity.AbsoluteExpiresAt.Valid {
		session.SetAbsoluteExpiresAt(entity.AbsoluteExpiresAt.Time)
	}

	// RotatedAt: pointer remains nil if NULL (session never refreshed)
	if entity.RotatedAt.Valid {
		rotated := entity.RotatedAt.Time
		session.SetRotatedAt(&rotated)
	}

	// UserAgent: empty string if NULL (client didn't provide user-agent)
	if entity.UserAgent.Valid {
		session.SetUserAgent(entity.UserAgent.String)
	}

	// IP: empty string if NULL (IP not captured)
	if entity.IP.Valid {
		session.SetIP(entity.IP.String)
	}

	// DeviceID: empty string if NULL (legacy session or device ID not provided)
	if entity.DeviceID.Valid {
		session.SetDeviceID(entity.DeviceID.String)
	}

	// LastRefreshAt: pointer remains nil if NULL (session never refreshed)
	if entity.LastRefreshAt.Valid {
		last := entity.LastRefreshAt.Time
		session.SetLastRefreshAt(&last)
	}

	return session
}
