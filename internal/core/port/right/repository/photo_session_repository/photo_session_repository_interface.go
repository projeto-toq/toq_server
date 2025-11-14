package photosessionrepository

import (
	"context"
	"database/sql"
	"time"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
)

// PhotoSessionRepositoryInterface defines persistence operations for photographer agenda entries and bookings.
type PhotoSessionRepositoryInterface interface {
	CreateEntries(ctx context.Context, tx *sql.Tx, entries []photosessionmodel.AgendaEntryInterface) ([]uint64, error)
	ListEntriesByRange(ctx context.Context, tx *sql.Tx, photographerID uint64, rangeStart, rangeEnd time.Time, entryType *photosessionmodel.AgendaEntryType) ([]photosessionmodel.AgendaEntryInterface, error)
	ListPhotographerIDs(ctx context.Context, tx *sql.Tx) ([]uint64, error)
	ListPhotographerIDsByLocation(ctx context.Context, tx *sql.Tx, city string, state string) ([]uint64, error)
	FindBlockingEntries(ctx context.Context, tx *sql.Tx, photographerID uint64, rangeStart, rangeEnd time.Time) ([]photosessionmodel.AgendaEntryInterface, error)
	DeleteEntriesBySource(ctx context.Context, tx *sql.Tx, photographerID uint64, entryType photosessionmodel.AgendaEntryType, source photosessionmodel.AgendaEntrySource, sourceID *uint64) (int64, error)
	GetEntryByID(ctx context.Context, tx *sql.Tx, entryID uint64) (photosessionmodel.AgendaEntryInterface, error)
	GetEntryByIDForUpdate(ctx context.Context, tx *sql.Tx, entryID uint64) (photosessionmodel.AgendaEntryInterface, error)
	DeleteEntryByID(ctx context.Context, tx *sql.Tx, entryID uint64) error
	UpdateEntrySourceID(ctx context.Context, tx *sql.Tx, entryID uint64, sourceID uint64) error
	UpdateEntry(ctx context.Context, tx *sql.Tx, entry photosessionmodel.AgendaEntryInterface) error

	CreateBooking(ctx context.Context, tx *sql.Tx, booking photosessionmodel.PhotoSessionBookingInterface) (uint64, error)
	UpdateBooking(ctx context.Context, tx *sql.Tx, booking photosessionmodel.PhotoSessionBookingInterface) error
	UpdateBookingStatus(ctx context.Context, tx *sql.Tx, bookingID uint64, status photosessionmodel.BookingStatus) error
	GetBookingByID(ctx context.Context, tx *sql.Tx, bookingID uint64) (photosessionmodel.PhotoSessionBookingInterface, error)
	GetBookingByIDForUpdate(ctx context.Context, tx *sql.Tx, bookingID uint64) (photosessionmodel.PhotoSessionBookingInterface, error)
	FindBookingByAgendaEntry(ctx context.Context, tx *sql.Tx, agendaEntryID uint64) (photosessionmodel.PhotoSessionBookingInterface, error)
	GetActiveBookingByListingIdentityID(ctx context.Context, tx *sql.Tx, listingIdentityID int64) (photosessionmodel.PhotoSessionBookingInterface, error)

	ListServiceAreasByPhotographer(ctx context.Context, tx *sql.Tx, photographerID uint64) ([]photosessionmodel.PhotographerServiceAreaInterface, error)
	GetServiceAreaByID(ctx context.Context, tx *sql.Tx, areaID uint64) (photosessionmodel.PhotographerServiceAreaInterface, error)
	ListAllServiceAreas(ctx context.Context, tx *sql.Tx, filter photosessionmodel.ServiceAreaFilter) ([]photosessionmodel.PhotographerServiceAreaInterface, error)
	CreateServiceArea(ctx context.Context, tx *sql.Tx, area photosessionmodel.PhotographerServiceAreaInterface) (uint64, error)
	UpdateServiceArea(ctx context.Context, tx *sql.Tx, area photosessionmodel.PhotographerServiceAreaInterface) error
	DeleteServiceArea(ctx context.Context, tx *sql.Tx, areaID uint64) error
}
