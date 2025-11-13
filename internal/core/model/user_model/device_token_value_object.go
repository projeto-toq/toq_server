package usermodel

// DeviceToken represents a push notification token associated with a user device
//
// This is a Value Object (immutable) used to track FCM/APNs tokens for sending
// push notifications to user devices. Device tokens are scoped to a user and
// can be associated with a specific device ID (UUIDv4).
//
// Design:
//   - Value Object (not an interface): simple data holder without behavior
//   - Immutable after creation: use setters on User to manage tokens
//   - No identity: equality based on token+deviceID combination
//
// Schema Mapping:
//   - Database table: device_tokens
//   - Foreign Key: user_id → users.id (CASCADE on delete)
//
// Use Cases:
//   - Push notifications: retrieve tokens for opted-in users
//   - Session management: associate tokens with device IDs
//   - Cleanup: remove tokens when user logs out or opts out
//
// Example:
//
//	token := DeviceToken{
//	    ID:       1,
//	    UserID:   123,
//	    Token:    "fcm_token_abc123...",
//	    DeviceID: "550e8400-e29b-41d4-a716-446655440000",
//	    Platform: "android",
//	}
type DeviceToken struct {
	// ID is the unique identifier of the device token record (AUTO_INCREMENT)
	ID int64

	// UserID is the foreign key to users.id (NOT NULL)
	// Associates this token with a specific user account
	UserID int64

	// Token is the FCM or APNs push notification token (NOT NULL, VARCHAR(255))
	// Format depends on push service (FCM: ~152 chars, APNs: 64 hex chars)
	// Example: "eHkN3-K8R9C_5LqQ2..."
	Token string

	// DeviceID is the unique device identifier (UUIDv4) (NULL, VARCHAR(100))
	// Used to associate token with specific device for session management
	// Example: "550e8400-e29b-41d4-a716-446655440000"
	DeviceID string

	// Platform is the device platform (NULL, VARCHAR(45))
	// Allowed values: "android", "ios", "web"
	// Used for platform-specific notification customization
	Platform string
}
