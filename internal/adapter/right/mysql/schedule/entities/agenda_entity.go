// Package scheduleentity holds the MySQL persistence shapes for schedule data (listing_agendas, listing_agenda_rules, listing_agenda_entries).
// These structs mirror the DB schema defined in scripts/db_creation.sql and are only used inside the MySQL adapter layer.
package scheduleentity

// AgendaEntity maps listing_agendas rows (InnoDB, utf8mb4) and is restricted to the MySQL adapter layer.
// Use scheduleconverters for domain conversions; do not import domain packages here.
type AgendaEntity struct {
	// ID is the AUTO_INCREMENT primary key (INT UNSIGNED NOT NULL)
	ID uint64
	// ListingIdentityID references listing_identities.id (INT UNSIGNED NOT NULL)
	ListingIdentityID int64
	// OwnerID references users.id (INT UNSIGNED NOT NULL)
	OwnerID int64
	// Timezone is the Olson timezone string (VARCHAR NOT NULL, e.g., "America/Sao_Paulo")
	Timezone string
}
