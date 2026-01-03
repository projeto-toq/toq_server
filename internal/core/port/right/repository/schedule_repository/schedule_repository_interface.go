// Package schedulerepository exposes the ScheduleRepositoryInterface contract mirrored by the MySQL adapter.
// Documentation rules (Section 8):
// - Describe purpose, transactional expectations, and returned infra errors (sql.ErrNoRows, fmt.Errorf).
// - Repositories never embed HTTP concerns; services map infra errors to business responses.
// - Callers manage transaction lifecycle; methods honor provided tx when not nil.
package schedulerepository

import (
	"context"
	"database/sql"
	"time"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
)

// ScheduleRepositoryInterface defines persistence operations for listing agendas and their rules/entries.
// Methods return pure infrastructure errors (sql.ErrNoRows, fmt.Errorf) and rely on services for business mapping.
type ScheduleRepositoryInterface interface {
	// GetAgendaByListingIdentityID returns the agenda for a listing_identity; sql.ErrNoRows when absent.
	GetAgendaByListingIdentityID(ctx context.Context, tx *sql.Tx, listingIdentityID int64) (schedulemodel.AgendaInterface, error)
	// InsertAgenda creates a new agenda and returns its ID; tx required for atomicity with related inserts.
	InsertAgenda(ctx context.Context, tx *sql.Tx, agenda schedulemodel.AgendaInterface) (uint64, error)
	// InsertRules bulk-inserts agenda rules; tx required; returns first infra error; sets rule IDs on success.
	InsertRules(ctx context.Context, tx *sql.Tx, rules []schedulemodel.AgendaRuleInterface) error
	// DeleteRulesByAgenda hard-deletes rules of an agenda; returns sql.ErrNoRows when nothing was removed.
	DeleteRulesByAgenda(ctx context.Context, tx *sql.Tx, agendaID uint64) error
	// ListRulesByAgenda lists all rules for the agenda ordered by weekday/start; empty slice when none.
	ListRulesByAgenda(ctx context.Context, tx *sql.Tx, agendaID uint64) ([]schedulemodel.AgendaRuleInterface, error)
	// GetRuleByID fetches a rule by id; sql.ErrNoRows when missing.
	GetRuleByID(ctx context.Context, tx *sql.Tx, ruleID uint64) (schedulemodel.AgendaRuleInterface, error)
	// UpdateRule updates a rule row; returns sql.ErrNoRows when the target no longer exists.
	UpdateRule(ctx context.Context, tx *sql.Tx, rule schedulemodel.AgendaRuleInterface) error
	// DeleteRule removes a single rule; returns sql.ErrNoRows when not found.
	DeleteRule(ctx context.Context, tx *sql.Tx, ruleID uint64) error
	// ListBlockRules lists blocking rules filtered by owner/listing/weekday; empty slice when none.
	ListBlockRules(ctx context.Context, tx *sql.Tx, filter schedulemodel.BlockRulesFilter) ([]schedulemodel.AgendaRuleInterface, error)
	// ListOwnerSummary returns consolidated summary per listing for an owner with pagination.
	ListOwnerSummary(ctx context.Context, tx *sql.Tx, filter schedulemodel.OwnerSummaryFilter) (schedulemodel.OwnerSummaryResult, error)
	// ListAgendaEntries returns paginated entries for an agenda range; empty slice when none.
	ListAgendaEntries(ctx context.Context, tx *sql.Tx, filter schedulemodel.AgendaDetailFilter) (schedulemodel.AgendaEntriesPage, error)
	// GetEntryByID fetches an entry by id; sql.ErrNoRows when missing.
	GetEntryByID(ctx context.Context, tx *sql.Tx, entryID uint64) (schedulemodel.AgendaEntryInterface, error)
	// GetEntryByVisitID fetches the entry linked to a visit; sql.ErrNoRows when absent.
	GetEntryByVisitID(ctx context.Context, tx *sql.Tx, visitID uint64) (schedulemodel.AgendaEntryInterface, error)
	// InsertEntry creates an entry and returns its ID; tx required; sets ID on the domain object.
	InsertEntry(ctx context.Context, tx *sql.Tx, entry schedulemodel.AgendaEntryInterface) (uint64, error)
	// UpdateEntry updates mutable fields of an entry; returns sql.ErrNoRows when target not found.
	UpdateEntry(ctx context.Context, tx *sql.Tx, entry schedulemodel.AgendaEntryInterface) error
	// DeleteEntry removes an entry; returns sql.ErrNoRows when nothing deleted.
	DeleteEntry(ctx context.Context, tx *sql.Tx, entryID uint64) error
	// ListEntriesBetween returns entries overlapping [from,to); empty slice when none.
	ListEntriesBetween(ctx context.Context, tx *sql.Tx, agendaID uint64, from time.Time, to time.Time) ([]schedulemodel.AgendaEntryInterface, error)
	// GetAvailabilityData aggregates rules and entries for availability computation; bubbles sql.ErrNoRows for missing agenda.
	GetAvailabilityData(ctx context.Context, tx *sql.Tx, filter schedulemodel.AvailabilityFilter) (schedulemodel.AvailabilityData, error)
}
