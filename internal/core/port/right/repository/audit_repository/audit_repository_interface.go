package auditrepository

import (
	"context"
	"database/sql"

	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
)

// Repository persists immutable audit events.
// Transactions are passed by the caller when batching with domain writes.
type Repository interface {
	CreateEvent(ctx context.Context, tx *sql.Tx, event auditmodel.AuditEvent) error
}
