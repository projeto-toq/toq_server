package userentity

import "database/sql"

// UserValidationEntity represents a row from the temp_user_validations table
//
// This entity stores temporary validation codes for email, phone, and password recovery.
// Used during account verification flows and profile update confirmations.
//
// Schema Mapping:
//   - Database: temp_user_validations table (InnoDB, utf8mb4_unicode_ci)
//   - Primary Key: user_id (FOREIGN KEY to users.id ON DELETE CASCADE)
//   - No AUTO_INCREMENT: user_id is set explicitly
//   - All code fields are nullable: NULL indicates no pending validation for that type
//
// Table Purpose:
//   - Email Change: NewEmail + EmailCode + EmailCodeExp
//   - Phone Change: NewPhone + PhoneCode + PhoneCodeExp
//   - Password Reset: PasswordCode + PasswordCodeExp (no "new value" stored for security)
//
// NULL Handling:
//   - sql.NullString: Used for VARCHAR columns (codes and new contact values)
//   - sql.NullTime: Used for DATETIME columns (expiration timestamps)
//   - All fields nullable to support partial validation states
//
// Lifecycle:
//   - Created when user requests email/phone change or password reset
//   - Codes expire after configured TTL (typically 5-10 minutes)
//   - Row deleted after successful validation or when all codes expired
//   - CASCADE DELETE: Removed automatically when user is deleted
//
// Conversion:
//   - To Domain: Use userconverters.ValidationEntityToDomain()
//   - From Domain: Use userconverters.ValidationDomainToEntity()
//
// Security Considerations:
//   - Codes are hashed before storage (enforced by service layer)
//   - Expiration timestamps prevent replay attacks
//   - No sensitive data stored (password hash remains in users.password)
//
// Important:
//   - DO NOT use this struct outside the adapter layer
//   - DO NOT add business logic methods to this struct
//   - DO NOT import core/model packages here
type UserValidationEntity struct {
	// UserID is the user's unique identifier (PRIMARY KEY, FOREIGN KEY to users.id)
	// ON DELETE CASCADE: validation row removed when user deleted
	UserID uint32

	// NewEmail stores the pending new email address (NULL, VARCHAR(45))
	// Set when user requests email change, cleared after successful validation
	// Example: "newemail@example.com"
	NewEmail sql.NullString

	// EmailCode is the hashed validation code sent to NewEmail (NULL, VARCHAR(6))
	// User must provide this code to confirm email ownership
	// Hashed with bcrypt before storage (service layer responsibility)
	EmailCode sql.NullString

	// EmailCodeExp is the expiration timestamp for EmailCode (NULL, TIMESTAMP(6))
	// After this time, code is invalid and user must request new code
	// Typically NOW() + 5 minutes
	EmailCodeExp sql.NullTime

	// NewPhone stores the pending new phone number in E.164 format (NULL, VARCHAR(25))
	// Set when user requests phone change, cleared after successful validation
	// Example: "+5511999998888"
	NewPhone sql.NullString

	// PhoneCode is the hashed validation code sent via SMS to NewPhone (NULL, VARCHAR(6))
	// User must provide this code to confirm phone ownership
	// Hashed with bcrypt before storage
	PhoneCode sql.NullString

	// PhoneCodeExp is the expiration timestamp for PhoneCode (NULL, TIMESTAMP(6))
	// After this time, code is invalid
	PhoneCodeExp sql.NullTime

	// PasswordCode is the hashed recovery code sent to user's email (NULL, VARCHAR(6))
	// Used for password reset flow (no "new password" stored here for security)
	// Hashed with bcrypt before storage
	PasswordCode sql.NullString

	// PasswordCodeExp is the expiration timestamp for PasswordCode (NULL, TIMESTAMP(6))
	// After this time, code is invalid
	PasswordCodeExp sql.NullTime
}
