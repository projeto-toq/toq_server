package photosessionservices

import (
	"context"
	"database/sql"
	"time"

	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	photosessionrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/photo_session_repository"
	userrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/user_repository"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	holidayservices "github.com/projeto-toq/toq_server/internal/core/service/holiday_service"
)

// PhotoSessionServiceInterface exposes orchestration helpers around photographer agenda entries.
type PhotoSessionServiceInterface interface {
	EnsurePhotographerAgendaWithTx(ctx context.Context, tx *sql.Tx, input EnsureAgendaInput) error
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
	ListServiceAreas(ctx context.Context, input ListServiceAreasInput) (ListServiceAreasOutput, error)
	CreateServiceArea(ctx context.Context, input CreateServiceAreaInput) (ServiceAreaResult, error)
	GetServiceArea(ctx context.Context, input ServiceAreaDetailInput) (ServiceAreaResult, error)
	UpdateServiceArea(ctx context.Context, input UpdateServiceAreaInput) (ServiceAreaResult, error)
	DeleteServiceArea(ctx context.Context, input DeleteServiceAreaInput) error
}

type photoSessionService struct {
	repo           photosessionrepository.PhotoSessionRepositoryInterface
	listingRepo    listingrepository.ListingRepoPortInterface
	userRepo       userrepository.UserRepoPortInterface
	holidayService holidayservices.HolidayServiceInterface
	globalService  globalservice.GlobalServiceInterface
	cfg            Config
	now            func() time.Time
}

// NewPhotoSessionService wires a photo session service with explicit config.
func NewPhotoSessionService(
	repo photosessionrepository.PhotoSessionRepositoryInterface,
	listingRepo listingrepository.ListingRepoPortInterface,
	userRepo userrepository.UserRepoPortInterface,
	holidayService holidayservices.HolidayServiceInterface,
	globalService globalservice.GlobalServiceInterface,
	cfg Config,
) PhotoSessionServiceInterface {
	return &photoSessionService{
		repo:           repo,
		listingRepo:    listingRepo,
		userRepo:       userRepo,
		holidayService: holidayService,
		globalService:  globalService,
		cfg:            cfg,
		now:            time.Now,
	}
}
