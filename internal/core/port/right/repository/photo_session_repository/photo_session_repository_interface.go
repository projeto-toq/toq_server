package photosessionrepository

import (
	"context"
	"database/sql"
	"time"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
)

// PhotoSessionRepositoryInterface defines persistence operations for photographer slots and bookings.
type PhotoSessionRepositoryInterface interface {
	ListAvailableSlots(ctx context.Context, tx *sql.Tx, params photosessionmodel.SlotListParams) ([]photosessionmodel.PhotographerSlotInterface, int64, error)
	GetSlotForUpdate(ctx context.Context, tx *sql.Tx, slotID uint64) (photosessionmodel.PhotographerSlotInterface, error)
	MarkSlotReserved(ctx context.Context, tx *sql.Tx, slotID uint64, token string, reservedUntil time.Time) error
	MarkSlotBooked(ctx context.Context, tx *sql.Tx, slotID uint64, bookedAt time.Time) error
	MarkSlotAvailable(ctx context.Context, tx *sql.Tx, slotID uint64) error
	InsertBooking(ctx context.Context, tx *sql.Tx, booking photosessionmodel.PhotoSessionBookingInterface) (uint64, error)
	GetBookingForUpdate(ctx context.Context, tx *sql.Tx, bookingID uint64) (photosessionmodel.PhotoSessionBookingInterface, error)
	UpdateBookingStatus(ctx context.Context, tx *sql.Tx, bookingID uint64, status photosessionmodel.BookingStatus) error
	BulkUpsertSlots(ctx context.Context, tx *sql.Tx, slots []photosessionmodel.PhotographerSlotInterface) error
	DeleteSlotsOutsideRange(ctx context.Context, tx *sql.Tx, photographerID uint64, windowStart, windowEnd time.Time) (int64, error)
	ListSlotsByRange(ctx context.Context, tx *sql.Tx, photographerID uint64, rangeStart, rangeEnd time.Time) ([]photosessionmodel.PhotographerSlotInterface, error)
	ListSlotsByRangePaginated(ctx context.Context, tx *sql.Tx, photographerID uint64, rangeStart, rangeEnd time.Time, limit, offset int) ([]photosessionmodel.PhotographerSlotInterface, int64, error)
	ListSlotsForPeriod(ctx context.Context, tx *sql.Tx, rangeStart, rangeEnd time.Time) ([]photosessionmodel.PhotographerSlotInterface, error)
	CreateTimeOff(ctx context.Context, tx *sql.Tx, timeOff photosessionmodel.PhotographerTimeOffInterface) (uint64, error)
	DeleteTimeOff(ctx context.Context, tx *sql.Tx, timeOffID uint64) error
	ListTimeOff(ctx context.Context, tx *sql.Tx, photographerID uint64, rangeStart, rangeEnd time.Time) ([]photosessionmodel.PhotographerTimeOffInterface, error)
	ListDefaultAvailability(ctx context.Context, tx *sql.Tx, photographerID uint64) ([]photosessionmodel.PhotographerDefaultAvailabilityInterface, error)
	ReplaceDefaultAvailability(ctx context.Context, tx *sql.Tx, photographerID uint64, records []photosessionmodel.PhotographerDefaultAvailabilityInterface) error
	DeleteDefaultAvailability(ctx context.Context, tx *sql.Tx, photographerID uint64, weekday *time.Weekday, period *photosessionmodel.SlotPeriod) error
}
