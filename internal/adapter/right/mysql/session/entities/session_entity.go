package sessionentities

import (
	"database/sql"
	"time"
)

// SessionEntity represents a row from the sessions table in the database
//
// This struct maps directly to the database schema and uses sql.Null* types
// for nullable columns. It should ONLY be used within the MySQL adapter layer.
//
// Schema Mapping:
//   - Database: sessions table (InnoDB, utf8mb4)
//   - Primary Key: id (INT UNSIGNED AUTO_INCREMENT)
//   - Unique Constraint: refresh_hash (CHAR(64) UNIQUE)
//   - Foreign Key: user_id → users(id) ON DELETE CASCADE
//   - Indexes: idx_sessions_user_active, idx_sessions_expires_at, idx_sessions_revoked, idx_sessions_token_jti
//
// NULL Handling:
//   - sql.NullString: Used for VARCHAR/CHAR columns that allow NULL (token_jti, user_agent, ip, device_id)
//   - sql.NullTime: Used for DATETIME(6) columns that allow NULL (absolute_expires_at, rotated_at, last_refresh_at)
//   - Direct types: Used for NOT NULL columns (id, user_id, refresh_hash, expires_at, created_at, revoked, rotation_counter)
//
// Conversion:
//   - To Domain: Use sessionconverters.SessionEntityToDomain()
//   - From Domain: Use sessionconverters.SessionDomainToEntity()
//
// Important:
//   - DO NOT use this struct outside the adapter layer
//   - DO NOT add business logic methods to this struct
//   - DO NOT import core/model packages here
type SessionEntity struct {
	// ID is the session's unique identifier (PRIMARY KEY, AUTO_INCREMENT)
	// Type: INT UNSIGNED, Range: 0 to 4,294,967,295
	ID int64

	// UserID is the ID of the user who owns this session (NOT NULL, FOREIGN KEY)
	// References users(id) with ON DELETE CASCADE
	// Type: INT UNSIGNED
	UserID int64

	// RefreshHash is the SHA-256 hash of the refresh token (NOT NULL, UNIQUE)
	// Used for refresh token validation and session lookup
	// Type: CHAR(64), Format: 64 hexadecimal characters
	RefreshHash string

	// TokenJTI is the unique identifier (JTI claim) of the JWT access token (NULL)
	// Used for token revocation and tracking
	// Type: CHAR(36), Format: UUIDv4 (e.g., "550e8400-e29b-41d4-a716-446655440000")
	// NULL when session created but no access token issued yet
	TokenJTI sql.NullString

	// ExpiresAt is the sliding expiration timestamp of the refresh token (NOT NULL)
	// Extended on each token refresh up to AbsoluteExpiresAt
	// Type: DATETIME(6), Precision: microseconds
	ExpiresAt time.Time

	// AbsoluteExpiresAt is the hard limit expiration timestamp (NULL)
	// Session cannot be refreshed after this timestamp regardless of activity
	// Type: DATETIME(6), Precision: microseconds
	// NULL means no absolute limit (session can be refreshed indefinitely)
	AbsoluteExpiresAt sql.NullTime

	// CreatedAt is the timestamp when the session was created (NOT NULL, DEFAULT CURRENT_TIMESTAMP(6))
	// Type: DATETIME(6), Precision: microseconds
	CreatedAt time.Time

	// RotatedAt is the timestamp of the last token rotation (NULL)
	// Updated when refresh token is rotated during token refresh flow
	// Type: DATETIME(6), Precision: microseconds
	// NULL for sessions that have never been refreshed
	RotatedAt sql.NullTime

	// UserAgent is the HTTP User-Agent header from the client (NULL)
	// Used for security tracking and device fingerprinting
	// Type: VARCHAR(255)
	// NULL when user-agent not provided by client
	UserAgent sql.NullString

	// IP is the client's IP address (NULL)
	// Used for security tracking and anomaly detection
	// Type: VARCHAR(64), Format: IPv4 (e.g., "192.168.1.1") or IPv6 (e.g., "2001:0db8::1")
	// NULL when IP cannot be determined
	IP sql.NullString

	// DeviceID is the unique device identifier from X-Device-Id header (NULL)
	// Type: VARCHAR(100), Format: UUIDv4
	// Used for multi-device session management
	// NULL when device ID not provided (legacy sessions)
	DeviceID sql.NullString

	// RotationCounter tracks the number of times the refresh token has been rotated (NOT NULL, DEFAULT 0)
	// Incremented on each token refresh
	// Type: INT UNSIGNED, Range: 0 to 4,294,967,295
	// Used for detecting token reuse attacks
	RotationCounter int

	// LastRefreshAt is the timestamp of the last successful token refresh (NULL)
	// Type: DATETIME(6), Precision: microseconds
	// NULL for sessions that have never been refreshed (only initial sign-in)
	LastRefreshAt sql.NullTime

	// Revoked indicates if the session has been revoked (NOT NULL, DEFAULT 0)
	// Type: TINYINT UNSIGNED, Values: 0 (active) or 1 (revoked)
	// Revoked sessions cannot be used for authentication
	// Note: Using bool for semantic clarity; DB stores 0/1
	Revoked bool
}
