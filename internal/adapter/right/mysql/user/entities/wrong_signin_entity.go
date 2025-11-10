package userentity

import "time"

// WrongSignInEntity represents a row from the temp_wrong_signin table
//
// This entity tracks failed sign-in attempts for security monitoring and brute-force prevention.
// Used to implement account lockout policies after repeated failed authentication attempts.
//
// Schema Mapping:
//   - Database: temp_wrong_signin table (InnoDB, utf8mb4_unicode_ci)
//   - Primary Key: user_id (FOREIGN KEY to users.id ON DELETE CASCADE)
//   - No AUTO_INCREMENT: user_id is set explicitly
//   - No nullable columns: all fields are required when tracking attempts
//
// Table Purpose:
//   - Count consecutive failed sign-in attempts
//   - Track timestamp of last failed attempt
//   - Support temporary account lockout after threshold reached
//   - Prevent brute-force password attacks
//
// Lifecycle:
//   - Created on first failed sign-in attempt
//   - Updated with each subsequent failure (increment counter, update timestamp)
//   - Deleted on successful sign-in (reset tracking)
//   - May be deleted by cleanup job after configured timeout (e.g., 24 hours)
//   - CASCADE DELETE: Removed automatically when user is deleted
//
// Conversion:
//   - To Domain: Use userconverters.WrongSigninEntityToDomain()
//   - From Domain: Use userconverters.WrongSigninDomainToEntity()
//
// Security Rules (enforced by service layer):
//   - Threshold typically 5 attempts within 15 minutes
//   - Account locked temporarily (e.g., 30 minutes) after threshold
//   - Counter resets on successful sign-in
//   - LastAttemptAT used to calculate lockout expiration
//
// Important:
//   - DO NOT use this struct outside the adapter layer
//   - DO NOT add business logic methods to this struct
//   - DO NOT import core/model packages here
type WrongSignInEntity struct {
	// UserID is the user's unique identifier (PRIMARY KEY, FOREIGN KEY to users.id, INT UNSIGNED)
	// ON DELETE CASCADE: tracking row removed when user deleted
	UserID uint32

	// FailedAttempts is the count of consecutive failed sign-in attempts (NOT NULL, TINYINT UNSIGNED)
	// Incremented on each failed authentication
	// Reset to 0 on successful sign-in
	// Max value: 255 (TINYINT limit, but lockout typically triggers at 5-10)
	FailedAttempts uint8

	// LastAttemptAT is the timestamp of the most recent failed sign-in (NOT NULL, TIMESTAMP(6))
	// Used to calculate lockout duration and attempt rate
	// Updated on each failed attempt
	// Example: used to check if 5 attempts happened within 15 minutes
	LastAttemptAT time.Time
}
