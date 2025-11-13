package scheduleservices

import (
	"context"
	"database/sql"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	schedulerepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/schedule_repository"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
)

// ScheduleServiceInterface exposes operations to orchestrate listing agendas.
type ScheduleServiceInterface interface {
	CreateDefaultAgenda(ctx context.Context, input CreateDefaultAgendaInput) (schedulemodel.AgendaInterface, error)
	CreateDefaultAgendaWithTx(ctx context.Context, tx *sql.Tx, input CreateDefaultAgendaInput) (schedulemodel.AgendaInterface, error)
	GetAgendaByListingIdentityID(ctx context.Context, listingIdentityID int64) (schedulemodel.AgendaInterface, error)
	ListBlockRules(ctx context.Context, filter schedulemodel.BlockRulesFilter) (schedulemodel.RuleListResult, error)
	CreateRules(ctx context.Context, input CreateRuleInput) (RuleMutationResult, error)
	UpdateRule(ctx context.Context, input UpdateRuleInput) (schedulemodel.AgendaRuleInterface, error)
	DeleteRule(ctx context.Context, input DeleteRuleInput) error
	ListRules(ctx context.Context, listingIdentityID, ownerID int64) (schedulemodel.RuleListResult, error)
	ListOwnerSummary(ctx context.Context, filter schedulemodel.OwnerSummaryFilter) (schedulemodel.OwnerSummaryResult, error)
	ListAgendaEntries(ctx context.Context, filter schedulemodel.AgendaDetailFilter) (schedulemodel.AgendaDetailResult, error)
	CreateBlockEntry(ctx context.Context, input CreateBlockEntryInput) (schedulemodel.AgendaEntryInterface, error)
	UpdateBlockEntry(ctx context.Context, input UpdateBlockEntryInput) (schedulemodel.AgendaEntryInterface, error)
	DeleteBlockEntry(ctx context.Context, input DeleteEntryInput) error
	GetAvailability(ctx context.Context, filter schedulemodel.AvailabilityFilter) (AvailabilityResult, error)
	FinishListingAgenda(ctx context.Context, input FinishListingAgendaInput) error
}

type scheduleService struct {
	scheduleRepo           schedulerepository.ScheduleRepositoryInterface
	listingRepo            listingrepository.ListingRepoPortInterface
	globalService          globalservice.GlobalServiceInterface
	defaultBlockRuleRanges []RuleTimeRange
	config                 Config
}

// NewScheduleService builds a schedule service with configuration overrides.
func NewScheduleService(
	scheduleRepo schedulerepository.ScheduleRepositoryInterface,
	listingRepo listingrepository.ListingRepoPortInterface,
	globalService globalservice.GlobalServiceInterface,
	config Config,
) ScheduleServiceInterface {
	// return buildScheduleService(scheduleRepo, listingRepo, globalService, config)
	config = config.ensureDefaults()
	return &scheduleService{
		scheduleRepo:           scheduleRepo,
		listingRepo:            listingRepo,
		globalService:          globalService,
		defaultBlockRuleRanges: config.DefaultBlockRuleRanges,
		config:                 config,
	}
}
