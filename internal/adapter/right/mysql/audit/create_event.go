package mysqlauditadapter

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"

	auditconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/audit/converters"
	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const insertAuditEvent = `INSERT INTO audit_events
 (occurred_at, actor_id, actor_role, actor_device_id, actor_ip, actor_user_agent,
  target_type, target_id, target_version, operation, metadata, request_id, trace_id)
 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

// CreateEvent persists a single audit event into the audit_events table.
func (a *AuditAdapter) CreateEvent(ctx context.Context, tx *sql.Tx, event auditmodel.AuditEvent) error {
	ctx, spanEnd, _ := utils.GenerateTracer(ctx)
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity, err := auditconverters.EventDomainToEntity(ctx, event)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.audit.create.convert_error", slog.Any("err", err))
		return fmt.Errorf("convert audit event: %w", err)
	}

	metadataPayload, err := marshalMetadata(entity.Metadata)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.audit.create.metadata_marshal_error", slog.Any("err", err))
		return fmt.Errorf("marshal audit metadata: %w", err)
	}

	res, execErr := a.ExecContext(ctx, tx, "insert", insertAuditEvent,
		entity.OccurredAt,
		entity.ActorID,
		entity.ActorRole,
		entity.ActorDeviceID,
		entity.ActorIP,
		entity.ActorUserAgent,
		entity.TargetType,
		entity.TargetID,
		entity.TargetVersion,
		entity.Operation,
		metadataPayload,
		entity.RequestID,
		entity.TraceID,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.audit.create.exec_error", slog.Any("err", execErr))
		return fmt.Errorf("insert audit_event: %w", execErr)
	}

	id, lastIDErr := res.LastInsertId()
	if lastIDErr != nil {
		utils.SetSpanError(ctx, lastIDErr)
		logger.Error("mysql.audit.create.last_insert_id_error", slog.Any("err", lastIDErr))
		return fmt.Errorf("audit_event last insert id: %w", lastIDErr)
	}

	event.SetID(id)
	return nil
}

func marshalMetadata(metadata any) ([]byte, error) {
	if metadata == nil {
		return json.Marshal(struct{}{})
	}
	return json.Marshal(metadata)
}
