package userentity

// RealtorsAgencyEntity represents a row from the realtors_agency table in the database
//
// This entity establishes many-to-many relationships between agencies and realtors.
// Allows one agency to manage multiple realtors, and one realtor to work with multiple agencies.
//
// Schema Mapping:
//   - Database: realtors_agency table (InnoDB, utf8mb3)
//   - Primary Key: id (INT UNSIGNED AUTO_INCREMENT)
//   - Foreign Keys: agency_id → users.id (CASCADE DELETE), realtor_id → users.id
//   - Indexes: fk_agency_idx, fk_realtor_idx
//
// Table Purpose:
//   - Link realtors to their managing agencies
//   - Support multi-agency realtors (freelance realtors)
//   - Track agency-realtor associations for commission splitting
//   - Enable agency-level filtering of listings and reports
//
// Lifecycle:
//   - Created when realtor accepts agency invitation
//   - Created when admin manually assigns realtor to agency
//   - Deleted when agency removes realtor from team
//   - Deleted when realtor leaves agency
//   - CASCADE DELETE: association removed when agency is deleted
//   - Standard DELETE: orphaned realtor remains when realtor user is deleted
//
// Conversion:
//   - To Domain: Use userconverters.RealtorsAgencyEntityToDomain()
//   - From Domain: Use userconverters.RealtorsAgencyDomainToEntity()
//
// Business Rules (enforced by service layer):
//   - agency_id must reference a user with "agency" role
//   - realtor_id must reference a user with "realtor" role
//   - No duplicate associations (same agency-realtor pair)
//   - At least one active association required for commission tracking
//
// Important:
//   - DO NOT use this struct outside the adapter layer
//   - DO NOT add business logic methods to this struct
//   - DO NOT import core/model packages here
type RealtorsAgencyEntity struct {
	// ID is the association's unique identifier (PRIMARY KEY, AUTO_INCREMENT, INT UNSIGNED)
	ID uint32

	// AgencyID is the managing agency's user ID (NOT NULL, INT UNSIGNED, FOREIGN KEY to users.id)
	// CASCADE DELETE: association removed when agency is deleted
	// Must reference a user with role "agency"
	AgencyID uint32

	// RealtorID is the realtor's user ID (NOT NULL, INT UNSIGNED, FOREIGN KEY to users.id)
	// Standard DELETE behavior (no CASCADE): allows tracking orphaned associations
	// Must reference a user with role "realtor"
	RealtorID uint32
}
