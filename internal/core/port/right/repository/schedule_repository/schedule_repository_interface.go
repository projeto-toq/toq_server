// Package schedulerepository exposes the ScheduleRepositoryInterface contract mirrored by the MySQL adapter.
// Documentation rules (Section 8):
// - Describe purpose, transactional expectations, and returned infra errors (sql.ErrNoRows, fmt.Errorf).
// - Repositories never embed HTTP concerns; services map infra errors to business responses.
// - Callers manage transaction lifecycle; methods honor provided tx when not nil.
//
// Usage expectations:
// - All methods must be called with a request-scoped context enriched by the service layer for logging/tracing.
// - Transactions are optional where explicitly stated; when provided, the adapter MUST honor the given tx.
// - Errors are pure infrastructure errors; absence is expressed via sql.ErrNoRows so services can map to domain/HTTP.
package schedulerepository

import (
	"context"
	"database/sql"
	"time"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
)

// ScheduleRepositoryInterface defines persistence operations for listing agendas and their rules/entries.
// Methods return pure infrastructure errors (sql.ErrNoRows, fmt.Errorf) and rely on services for business mapping.
// Each method must initialize tracing (utils.GenerateTracer) inside the adapter and use InstrumentedAdapter for queries.
type ScheduleRepositoryInterface interface {
	// GetAgendaByListingIdentityID returns the agenda for a listing_identity_id.
	// Returns sql.ErrNoRows when the agenda does not exist; uses the provided transaction for consistent reads.
	GetAgendaByListingIdentityID(ctx context.Context, tx *sql.Tx, listingIdentityID int64) (schedulemodel.AgendaInterface, error)
	// InsertAgenda creates a new agenda row and returns its generated ID.
	// Must be executed inside a transaction when chained with rule/entry inserts to guarantee atomicity.
	InsertAgenda(ctx context.Context, tx *sql.Tx, agenda schedulemodel.AgendaInterface) (uint64, error)
	// InsertRules bulk-inserts agenda rules, setting generated IDs on the passed domain objects.
	// Returns the first infrastructure error; expects a non-nil transaction for atomic writes.
	InsertRules(ctx context.Context, tx *sql.Tx, rules []schedulemodel.AgendaRuleInterface) error
	// DeleteRulesByAgenda hard-deletes all rules of an agenda; returns sql.ErrNoRows when no row is removed.
	DeleteRulesByAgenda(ctx context.Context, tx *sql.Tx, agendaID uint64) error
	// ListRulesByAgenda lists all rules ordered by weekday and start minute; returns empty slice when none.
	ListRulesByAgenda(ctx context.Context, tx *sql.Tx, agendaID uint64) ([]schedulemodel.AgendaRuleInterface, error)
	// GetRuleByID fetches a single rule by id; returns sql.ErrNoRows when missing.
	GetRuleByID(ctx context.Context, tx *sql.Tx, ruleID uint64) (schedulemodel.AgendaRuleInterface, error)
	// UpdateRule updates a rule row; returns sql.ErrNoRows when the target no longer exists.
	UpdateRule(ctx context.Context, tx *sql.Tx, rule schedulemodel.AgendaRuleInterface) error
	// DeleteRule removes a single rule; returns sql.ErrNoRows when no row matches the given id.
	DeleteRule(ctx context.Context, tx *sql.Tx, ruleID uint64) error
	// ListBlockRules lists blocking rules filtered by owner, listing, and optional weekdays; returns empty slice when none.
	ListBlockRules(ctx context.Context, tx *sql.Tx, filter schedulemodel.BlockRulesFilter) ([]schedulemodel.AgendaRuleInterface, error)
	// ListOwnerSummary returns consolidated agenda summaries per listing for an owner with pagination.
	// Returns empty items when none found; total always reflects COUNT(*) for pagination.
	ListOwnerSummary(ctx context.Context, tx *sql.Tx, filter schedulemodel.OwnerSummaryFilter) (schedulemodel.OwnerSummaryResult, error)
	// ListAgendaEntries returns paginated entries for a listing agenda within a time window; empty slice when none.
	ListAgendaEntries(ctx context.Context, tx *sql.Tx, filter schedulemodel.AgendaDetailFilter) (schedulemodel.AgendaEntriesPage, error)
	// GetEntryByID fetches an agenda entry by id; returns sql.ErrNoRows when missing.
	GetEntryByID(ctx context.Context, tx *sql.Tx, entryID uint64) (schedulemodel.AgendaEntryInterface, error)
	// GetEntryByVisitID fetches the agenda entry linked to a visit id; returns sql.ErrNoRows when absent.
	GetEntryByVisitID(ctx context.Context, tx *sql.Tx, visitID uint64) (schedulemodel.AgendaEntryInterface, error)
	// InsertEntry creates a new agenda entry, sets its ID on the domain object, and returns the generated ID.
	// Requires an active transaction to maintain atomicity with related writes.
	InsertEntry(ctx context.Context, tx *sql.Tx, entry schedulemodel.AgendaEntryInterface) (uint64, error)
	// UpdateEntry updates mutable fields of an agenda entry; returns sql.ErrNoRows when the target does not exist.
	UpdateEntry(ctx context.Context, tx *sql.Tx, entry schedulemodel.AgendaEntryInterface) error
	// DeleteEntry removes an agenda entry; returns sql.ErrNoRows when nothing is deleted.
	DeleteEntry(ctx context.Context, tx *sql.Tx, entryID uint64) error
	// ListEntriesBetween returns entries overlapping [from, to); returns empty slice when none match.
	ListEntriesBetween(ctx context.Context, tx *sql.Tx, agendaID uint64, from time.Time, to time.Time) ([]schedulemodel.AgendaEntryInterface, error)
	// GetAvailabilityData aggregates rules and entries to compute availability for a listing.
	// Returns sql.ErrNoRows when the agenda is missing; otherwise bubbles infra errors for service mapping.
	GetAvailabilityData(ctx context.Context, tx *sql.Tx, filter schedulemodel.AvailabilityFilter) (schedulemodel.AvailabilityData, error)
}
