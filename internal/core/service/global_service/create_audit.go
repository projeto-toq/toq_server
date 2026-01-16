package globalservice

import (
	"context"
	"database/sql"
	"time"

	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateAudit persists an immutable audit trail record for data-changing operations.
// When executedBY is omitted the user ID is extracted from the request context.
func (gs *globalService) CreateAudit(ctx context.Context, tx *sql.Tx, table globalmodel.TableName, action string, executedBY ...int64) error {
	ctx, spanEnd, tracerErr := utils.GenerateTracer(ctx)
	if tracerErr != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("global.audit.tracer_error", "err", tracerErr)
		return utils.InternalError("Failed to initialize audit tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if gs.auditService == nil {
		logger.Error("global.audit.audit_service_nil")
		return utils.InternalError("Audit service not configured")
	}

	actorID := int64(0)
	if len(executedBY) > 0 {
		actorID = executedBY[0]
	} else if v := ctx.Value(globalmodel.TokenKey); v != nil {
		if infos, ok := v.(usermodel.UserInfos); ok {
			actorID = infos.ID
		}
	}

	record := auditmodel.RecordInput{
		Actor: auditmodel.AuditActor{ID: actorID},
		Target: auditmodel.AuditTarget{
			Type: mapTableToTarget(table),
			ID:   0, // legacy calls do not provide the specific entity ID
		},
		Operation: auditmodel.OperationUpdate,
		Metadata: map[string]any{
			"legacy_action": action,
			"table":         table.String(),
		},
		OccurredAt: time.Now().UTC(),
	}

	if err := gs.auditService.RecordChange(ctx, tx, record); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("global.audit.persist_error", "err", err, "table", table, "action", action)
		return err
	}

	return nil
}

func mapTableToTarget(table globalmodel.TableName) auditmodel.TargetType {
	switch table {
	case globalmodel.TableUsers:
		return auditmodel.TargetType("users")
	case globalmodel.TableListings:
		return auditmodel.TargetListingIdentity
	case globalmodel.TableProposals:
		return auditmodel.TargetProposal
	case globalmodel.TableRealtorAgency:
		return auditmodel.TargetType("realtors_agency")
	default:
		return auditmodel.TargetType(table.String())
	}
}
