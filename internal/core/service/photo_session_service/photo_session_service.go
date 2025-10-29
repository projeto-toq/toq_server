package photosessionservices

import (
	"context"
	"database/sql"
	"time"

	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	photosessionrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/photo_session_repository"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	holidayservices "github.com/projeto-toq/toq_server/internal/core/service/holiday_service"
)

// PhotoSessionServiceInterface exposes orchestration helpers around photographer agenda entries.
type PhotoSessionServiceInterface interface {
	EnsurePhotographerAgenda(ctx context.Context, input EnsureAgendaInput) error
	EnsurePhotographerAgendaWithTx(ctx context.Context, tx *sql.Tx, input EnsureAgendaInput) error
	RefreshPhotographerAgenda(ctx context.Context, input EnsureAgendaInput) error
	RefreshPhotographerAgendaWithTx(ctx context.Context, tx *sql.Tx, input EnsureAgendaInput) error
	CreateTimeOff(ctx context.Context, input TimeOffInput) (uint64, error)
	CreateTimeOffWithTx(ctx context.Context, tx *sql.Tx, input TimeOffInput) (uint64, error)
	DeleteTimeOff(ctx context.Context, input DeleteTimeOffInput) error
	DeleteTimeOffWithTx(ctx context.Context, tx *sql.Tx, input DeleteTimeOffInput) error
	ListTimeOff(ctx context.Context, input ListTimeOffInput) (ListTimeOffOutput, error)
	GetTimeOffDetail(ctx context.Context, input TimeOffDetailInput) (TimeOffDetailResult, error)
	UpdateTimeOff(ctx context.Context, input UpdateTimeOffInput) (TimeOffDetailResult, error)
	UpdateSessionStatus(ctx context.Context, input UpdateSessionStatusInput) error
	ListAgenda(ctx context.Context, input ListAgendaInput) (ListAgendaOutput, error)
	ListAvailability(ctx context.Context, input ListAvailabilityInput) (ListAvailabilityOutput, error)
	ReservePhotoSession(ctx context.Context, input ReserveSessionInput) (ReserveSessionOutput, error)
	ConfirmPhotoSession(ctx context.Context, input ConfirmSessionInput) (ConfirmSessionOutput, error)
	CancelPhotoSession(ctx context.Context, input CancelSessionInput) (CancelSessionOutput, error)
}

type photoSessionService struct {
	repo           photosessionrepository.PhotoSessionRepositoryInterface
	holidayRepo    photosessionrepository.PhotographerHolidayCalendarRepository
	listingRepo    listingrepository.ListingRepoPortInterface
	holidayService holidayservices.HolidayServiceInterface
	globalService  globalservice.GlobalServiceInterface
	cfg            Config
	now            func() time.Time
}

// NewPhotoSessionService wires a photo session service with its dependencies.
func NewPhotoSessionService(
	repo photosessionrepository.PhotoSessionRepositoryInterface,
	listingRepo listingrepository.ListingRepoPortInterface,
	holidayService holidayservices.HolidayServiceInterface,
	globalService globalservice.GlobalServiceInterface,
) PhotoSessionServiceInterface {
	return NewPhotoSessionServiceWithConfig(repo, listingRepo, holidayService, globalService, Config{})
}

// NewPhotoSessionServiceWithConfig wires a photo session service with explicit config.
func NewPhotoSessionServiceWithConfig(
	repo photosessionrepository.PhotoSessionRepositoryInterface,
	listingRepo listingrepository.ListingRepoPortInterface,
	holidayService holidayservices.HolidayServiceInterface,
	globalService globalservice.GlobalServiceInterface,
	cfg Config,
) PhotoSessionServiceInterface {
	var holidayRepo photosessionrepository.PhotographerHolidayCalendarRepository
	if value, ok := repo.(photosessionrepository.PhotographerHolidayCalendarRepository); ok {
		holidayRepo = value
	}

	return &photoSessionService{
		repo:           repo,
		holidayRepo:    holidayRepo,
		listingRepo:    listingRepo,
		holidayService: holidayService,
		globalService:  globalService,
		cfg:            normalizeConfig(cfg),
		now:            time.Now,
	}
}
