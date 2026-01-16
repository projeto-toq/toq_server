package auditservice

import (
	"context"
	"database/sql"

	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	auditrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/audit_repository"
)

// AuditServiceInterface defines operations to record audit events.
type AuditServiceInterface interface {
	RecordChange(ctx context.Context, tx *sql.Tx, input auditmodel.RecordInput) error
}

// auditService is the concrete implementation backed by a repository.
type auditService struct {
	repo auditrepository.Repository
}

// NewAuditService constructs a new audit service instance.
func NewAuditService(repo auditrepository.Repository) AuditServiceInterface {
	return &auditService{repo: repo}
}
