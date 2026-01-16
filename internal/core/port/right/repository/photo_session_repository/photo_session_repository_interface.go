// Package photosessionrepository declares the persistence contract for photographer agenda entries,
// bookings and service areas. Implementations must follow tracing, logging and transactional
// expectations defined in docs/toq_server_go_guide.md and return pure infrastructure errors
// (e.g., sql.ErrNoRows) without HTTP coupling.
package photosessionrepository

import (
	"context"
	"database/sql"
	"time"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
)

// PhotoSessionRepositoryInterface defines all persistence operations for photographer agenda entries,
// bookings and service areas. Methods accept a context for tracing/logging, an optional transaction
// (tx may be nil when atomicity is not required) and return sql.ErrNoRows in not-found scenarios.
// Concurrency-sensitive reads expose FOR UPDATE variants to allow callers to compose locking flows.
type PhotoSessionRepositoryInterface interface {
	// CreateEntries inserts multiple agenda entries; preserves input order in returned IDs.
	// Returns sql.ErrNoRows only when a non-nil transaction affects zero rows.
	CreateEntries(ctx context.Context, tx *sql.Tx, entries []photosessionmodel.AgendaEntryInterface) ([]uint64, error)
	// ListEntriesByRange returns agenda entries overlapping the provided window; when entryType is set, filters by type.
	// tx optional; returns empty slice when no entries match.
	ListEntriesByRange(ctx context.Context, tx *sql.Tx, photographerID uint64, rangeStart, rangeEnd time.Time, entryType *photosessionmodel.AgendaEntryType) ([]photosessionmodel.AgendaEntryInterface, error)
	// ListPhotographerIDs lists distinct photographer user IDs with active role; tx optional; empty slice if none.
	ListPhotographerIDs(ctx context.Context, tx *sql.Tx) ([]uint64, error)
	// ListPhotographerIDsByLocation lists photographer IDs that serve the given city/state; tx optional; empty slice if none.
	ListPhotographerIDsByLocation(ctx context.Context, tx *sql.Tx, city string, state string) ([]uint64, error)
	// FindBlockingEntries returns blocking agenda items overlapping the window; tx optional; empty slice if none.
	FindBlockingEntries(ctx context.Context, tx *sql.Tx, photographerID uint64, rangeStart, rangeEnd time.Time) ([]photosessionmodel.AgendaEntryInterface, error)
	// DeleteEntriesBySource deletes entries by source/source_id; tx required for atomicity; returns rows deleted (0 → sql.ErrNoRows).
	DeleteEntriesBySource(ctx context.Context, tx *sql.Tx, photographerID uint64, entryType photosessionmodel.AgendaEntryType, source photosessionmodel.AgendaEntrySource, sourceID *uint64) (int64, error)
	// GetEntryByID fetches an agenda entry; tx optional; returns sql.ErrNoRows when absent.
	GetEntryByID(ctx context.Context, tx *sql.Tx, entryID uint64) (photosessionmodel.AgendaEntryInterface, error)
	// GetEntryByIDForUpdate fetches an agenda entry with FOR UPDATE; tx must be non-nil for lock to take effect.
	GetEntryByIDForUpdate(ctx context.Context, tx *sql.Tx, entryID uint64) (photosessionmodel.AgendaEntryInterface, error)
	// DeleteEntryByID removes an agenda entry; tx optional; sql.ErrNoRows when id not found.
	DeleteEntryByID(ctx context.Context, tx *sql.Tx, entryID uint64) error
	// UpdateEntrySourceID sets source_id; tx optional; sql.ErrNoRows when id not found.
	UpdateEntrySourceID(ctx context.Context, tx *sql.Tx, entryID uint64, sourceID uint64) error
	// UpdateEntry overwrites entry fields; tx optional; sql.ErrNoRows when id not found.
	UpdateEntry(ctx context.Context, tx *sql.Tx, entry photosessionmodel.AgendaEntryInterface) error

	// CreateBooking creates a booking row linked to an agenda entry; tx required; returns new ID.
	CreateBooking(ctx context.Context, tx *sql.Tx, booking photosessionmodel.PhotoSessionBookingInterface) (uint64, error)
	// UpdateBooking updates booking fields (including times and reason); tx optional; sql.ErrNoRows when id not found.
	UpdateBooking(ctx context.Context, tx *sql.Tx, booking photosessionmodel.PhotoSessionBookingInterface) error
	// UpdateBookingStatus sets booking status; tx optional; sql.ErrNoRows when id not found.
	UpdateBookingStatus(ctx context.Context, tx *sql.Tx, bookingID uint64, status photosessionmodel.BookingStatus) error
	// GetBookingByID fetches a booking by id; tx optional; returns sql.ErrNoRows when absent.
	GetBookingByID(ctx context.Context, tx *sql.Tx, bookingID uint64) (photosessionmodel.PhotoSessionBookingInterface, error)
	// GetBookingByIDForUpdate fetches a booking with FOR UPDATE; tx must be non-nil for lock to apply.
	GetBookingByIDForUpdate(ctx context.Context, tx *sql.Tx, bookingID uint64) (photosessionmodel.PhotoSessionBookingInterface, error)
	// FindBookingByAgendaEntry fetches booking by agenda entry; tx optional; returns sql.ErrNoRows when absent.
	FindBookingByAgendaEntry(ctx context.Context, tx *sql.Tx, agendaEntryID uint64) (photosessionmodel.PhotoSessionBookingInterface, error)
	// GetActiveBookingByListingIdentityID returns latest active booking for a listing identity (statuses pending/accepted/active); sql.ErrNoRows when none.
	GetActiveBookingByListingIdentityID(ctx context.Context, tx *sql.Tx, listingIdentityID int64) (photosessionmodel.PhotoSessionBookingInterface, error)

	// ListServiceAreasByPhotographer lists service areas for a photographer; tx optional; empty slice when none.
	ListServiceAreasByPhotographer(ctx context.Context, tx *sql.Tx, photographerID uint64) ([]photosessionmodel.PhotographerServiceAreaInterface, error)
	// GetServiceAreaByID fetches a service area; tx optional; sql.ErrNoRows when absent.
	GetServiceAreaByID(ctx context.Context, tx *sql.Tx, areaID uint64) (photosessionmodel.PhotographerServiceAreaInterface, error)
	// ListAllServiceAreas lists service areas using the provided filter (city/state/limit/offset); tx optional; empty slice when none.
	ListAllServiceAreas(ctx context.Context, tx *sql.Tx, filter photosessionmodel.ServiceAreaFilter) ([]photosessionmodel.PhotographerServiceAreaInterface, error)
	// CreateServiceArea inserts a service area; tx required when coupled with other writes; returns new ID.
	CreateServiceArea(ctx context.Context, tx *sql.Tx, area photosessionmodel.PhotographerServiceAreaInterface) (uint64, error)
	// UpdateServiceArea updates city/state; tx optional; sql.ErrNoRows when id not found.
	UpdateServiceArea(ctx context.Context, tx *sql.Tx, area photosessionmodel.PhotographerServiceAreaInterface) error
	// DeleteServiceArea removes a service area; tx optional; sql.ErrNoRows when id not found.
	DeleteServiceArea(ctx context.Context, tx *sql.Tx, areaID uint64) error

	// DeleteOldBookings removes bookings in terminal statuses whose ends_at is older than cutoff; returns rows deleted.
	DeleteOldBookings(ctx context.Context, tx *sql.Tx, cutoff time.Time, limit int) (int64, error)
	// DeleteOldAgendaEntries removes agenda entries whose ends_at is older than cutoff and are not referenced by bookings; returns rows deleted.
	DeleteOldAgendaEntries(ctx context.Context, tx *sql.Tx, cutoff time.Time, limit int) (int64, error)
}
