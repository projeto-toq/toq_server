package userentity

import "database/sql"

// DeviceTokenEntity represents a row in the device_tokens table
//
// This entity maps directly to the database schema and uses sql.Null* types
// for nullable columns. It should ONLY be used within the MySQL adapter layer.
//
// Schema Mapping:
//   - Database table: device_tokens (InnoDB, utf8mb4)
//   - Primary Key: id (INT UNSIGNED AUTO_INCREMENT)
//   - Foreign Key: user_id → users.id (CASCADE on DELETE)
//   - Unique Constraint: (user_id, device_id) - ensures one token per device
//   - Indexes:
//   - PRIMARY KEY (id)
//   - FOREIGN KEY idx_device_tokens_user (user_id)
//   - UNIQUE KEY uk_device_tokens_user_device (user_id, device_id)
//
// NULL Handling:
//   - sql.NullString: Used for platform (VARCHAR(45) NULL)
//   - Direct types: Used for NOT NULL columns (id, user_id, device_token, device_id)
//
// Conversion:
//   - To Domain: Use userconverters.DeviceTokenEntityToVO()
//   - From Domain: Use userconverters.DeviceTokenVOToEntity()
//
// Important:
//   - DO NOT use this struct outside the adapter layer
//   - DO NOT add business logic methods to this struct
//   - DO NOT import core/model packages here
type DeviceTokenEntity struct {
	// ID is the device token record's unique identifier (PRIMARY KEY, AUTO_INCREMENT)
	// Generated automatically on INSERT
	ID int64

	// UserID is the foreign key to users.id (NOT NULL, INT UNSIGNED)
	// Associates this token with a specific user account
	// CASCADE DELETE: token is removed when user is deleted
	UserID int64

	// Token is the FCM or APNs push notification token (NOT NULL, VARCHAR(255))
	// Format varies by service:
	//   - FCM: ~152 characters (base64-like)
	//   - APNs: 64 characters (hex)
	// Example: "eHkN3-K8R9C_5LqQ2..."
	Token string

	// DeviceID is the unique device identifier (NOT NULL, VARCHAR(100))
	// Expected format: UUIDv4 (36 characters with hyphens)
	// Used for session management and device-specific operations
	// Example: "550e8400-e29b-41d4-a716-446655440000"
	// UNIQUE constraint with user_id: (user_id, device_id) must be unique
	DeviceID string

	// Platform is the device platform (NULL, VARCHAR(45))
	// Allowed values: "android", "ios", "web", NULL
	// Used for platform-specific notification customization
	// NULL indicates platform is unknown or not relevant
	Platform sql.NullString
}
