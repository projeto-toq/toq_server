package userentity

// AgencyInvite represents a row from the agency_invites table
//
// This entity stores pending invitations sent by agencies to realtors.
// Used to track and manage agency-realtor relationships before realtor accepts.
//
// Schema Mapping:
//   - Database: agency_invites table (InnoDB, utf8mb3)
//   - Primary Key: id (INT UNSIGNED AUTO_INCREMENT)
//   - Foreign Key: agency_id → users.id
//   - Index: fk_agency_invite_idx on agency_id
//
// Table Purpose:
//   - Track pending invitations from agencies to realtors
//   - Prevent duplicate invitations to same phone number
//   - Validate phone number before creating user account
//   - Support invitation acceptance flow
//
// Lifecycle:
//   - Created when agency sends invitation to realtor's phone
//   - Remains until realtor creates account and accepts
//   - Deleted after successful acceptance (relationship created in realtors_agency)
//   - May be deleted if agency cancels invitation
//   - May expire after configured timeout (e.g., 30 days)
//
// Conversion:
//   - To Domain: Use userconverters.AgencyInviteEntityToDomain()
//   - From Domain: Use userconverters.AgencyInviteDomainToEntity()
//
// Business Rules (enforced by service layer):
//   - One agency can invite same phone number only once (prevent spam)
//   - Phone number must be in E.164 format
//   - Agency must have active "agency" role
//   - Invitation deleted after realtor accepts
//
// Important:
//   - DO NOT use this struct outside the adapter layer
//   - DO NOT add business logic methods to this struct
//   - DO NOT import core/model packages here
type AgencyInvite struct {
	// ID is the invitation's unique identifier (PRIMARY KEY, AUTO_INCREMENT, INT UNSIGNED)
	ID int64

	// AgencyID is the inviting agency's user ID (NOT NULL, INT UNSIGNED, FOREIGN KEY to users.id)
	// References a user with "agency" role
	// Used to track which agency sent the invitation
	AgencyID int64

	// PhoneNumber is the invited realtor's phone in E.164 format (NOT NULL, VARCHAR(15))
	// Must be unique per agency (one agency cannot invite same number twice)
	// Used to match against user registration during acceptance flow
	// Example: "+5511999999999"
	PhoneNumber string
}
