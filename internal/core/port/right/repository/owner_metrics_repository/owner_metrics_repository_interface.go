package ownermetricsrepository

import (
	"context"
	"database/sql"
	"time"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// VisitResponseInput encapsulates the delta between request and the owner's first action for visits.
type VisitResponseInput struct {
	OwnerID      int64
	DeltaSeconds int64
	RespondedAt  time.Time
}

// ProposalResponseInput aggregates SLA data for proposal workflows.
type ProposalResponseInput struct {
	OwnerID      int64
	DeltaSeconds int64
	RespondedAt  time.Time
}

// Repository persists and retrieves aggregated SLA metrics per owner.
type Repository interface {
	UpsertVisitResponse(ctx context.Context, tx *sql.Tx, input VisitResponseInput) error
	UpsertProposalResponse(ctx context.Context, tx *sql.Tx, input ProposalResponseInput) error
	GetByOwnerID(ctx context.Context, tx *sql.Tx, ownerID int64) (usermodel.OwnerResponseMetrics, error)
}
