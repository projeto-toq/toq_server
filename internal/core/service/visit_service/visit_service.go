package visitservice

import (
	"context"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	schedulerepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/schedule_repository"
	visitrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/visit_repository"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	scheduleservices "github.com/projeto-toq/toq_server/internal/core/service/schedule_service"
)

// VisitDetailOutput aggregates the visit entity with the related listing snapshot.
type VisitDetailOutput struct {
	Visit   listingmodel.VisitInterface
	Listing listingmodel.ListingInterface
}

// VisitListOutput mirrors VisitDetailOutput with pagination metadata for list endpoints.
type VisitListOutput struct {
	Items []VisitDetailOutput
	Total int64
	Page  int
	Limit int
}

// Service exposes visit orchestration operations.
type Service interface {
	CreateVisit(ctx context.Context, input CreateVisitInput) (listingmodel.VisitInterface, error)
	ApproveVisit(ctx context.Context, visitID int64, ownerNotes string) (listingmodel.VisitInterface, error)
	RejectVisit(ctx context.Context, visitID int64, reason string) (listingmodel.VisitInterface, error)
	CancelVisit(ctx context.Context, visitID int64, reason string) (listingmodel.VisitInterface, error)
	CompleteVisit(ctx context.Context, visitID int64, ownerNotes string) (listingmodel.VisitInterface, error)
	MarkNoShow(ctx context.Context, visitID int64, ownerNotes string) (listingmodel.VisitInterface, error)
	GetVisit(ctx context.Context, visitID int64) (VisitDetailOutput, error)
	ListVisits(ctx context.Context, filter listingmodel.VisitListFilter) (VisitListOutput, error)
}

// NewService wires the visit service dependencies.
func NewService(gs globalservice.GlobalServiceInterface, visitRepo visitrepository.VisitRepositoryInterface, listingRepo listingrepository.ListingRepoPortInterface, scheduleRepo schedulerepository.ScheduleRepositoryInterface, scheduleSvc scheduleservices.ScheduleServiceInterface, config Config) Service {
	return &visitService{
		globalService: gs,
		visitRepo:     visitRepo,
		listingRepo:   listingRepo,
		scheduleRepo:  scheduleRepo,
		scheduleSvc:   scheduleSvc,
		config:        config,
	}
}

type visitService struct {
	globalService globalservice.GlobalServiceInterface
	visitRepo     visitrepository.VisitRepositoryInterface
	listingRepo   listingrepository.ListingRepoPortInterface
	scheduleRepo  schedulerepository.ScheduleRepositoryInterface
	scheduleSvc   scheduleservices.ScheduleServiceInterface
	config        Config
}
