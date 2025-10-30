package schedulerepository

import (
	"context"
	"database/sql"
	"time"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
)

// ScheduleRepositoryInterface defines persistence operations for listing agendas.
type ScheduleRepositoryInterface interface {
	GetAgendaByListingID(ctx context.Context, tx *sql.Tx, listingID int64) (schedulemodel.AgendaInterface, error)
	InsertAgenda(ctx context.Context, tx *sql.Tx, agenda schedulemodel.AgendaInterface) (uint64, error)
	InsertRules(ctx context.Context, tx *sql.Tx, rules []schedulemodel.AgendaRuleInterface) error
	DeleteRulesByAgenda(ctx context.Context, tx *sql.Tx, agendaID uint64) error
	ListRulesByAgenda(ctx context.Context, tx *sql.Tx, agendaID uint64) ([]schedulemodel.AgendaRuleInterface, error)
	GetRuleByID(ctx context.Context, tx *sql.Tx, ruleID uint64) (schedulemodel.AgendaRuleInterface, error)
	UpdateRule(ctx context.Context, tx *sql.Tx, rule schedulemodel.AgendaRuleInterface) error
	DeleteRule(ctx context.Context, tx *sql.Tx, ruleID uint64) error
	ListOwnerSummary(ctx context.Context, tx *sql.Tx, filter schedulemodel.OwnerSummaryFilter) (schedulemodel.OwnerSummaryResult, error)
	ListAgendaEntries(ctx context.Context, tx *sql.Tx, filter schedulemodel.AgendaDetailFilter) (schedulemodel.AgendaEntriesPage, error)
	GetEntryByID(ctx context.Context, tx *sql.Tx, entryID uint64) (schedulemodel.AgendaEntryInterface, error)
	InsertEntry(ctx context.Context, tx *sql.Tx, entry schedulemodel.AgendaEntryInterface) (uint64, error)
	UpdateEntry(ctx context.Context, tx *sql.Tx, entry schedulemodel.AgendaEntryInterface) error
	DeleteEntry(ctx context.Context, tx *sql.Tx, entryID uint64) error
	ListEntriesBetween(ctx context.Context, tx *sql.Tx, agendaID uint64, from time.Time, to time.Time) ([]schedulemodel.AgendaEntryInterface, error)
	GetAvailabilityData(ctx context.Context, tx *sql.Tx, filter schedulemodel.AvailabilityFilter) (schedulemodel.AvailabilityData, error)
}
