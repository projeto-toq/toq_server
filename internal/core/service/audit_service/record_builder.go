package auditservice

import (
	"context"

	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// BuildRecordFromContext assembles a RecordInput with actor and correlation metadata derived from the context.
//
// It populates:
//   - Actor: built from context (role/device/ip/user-agent) with the provided userID fallback.
//   - Target: uses the given target (defaults version to zero when nil to satisfy DB NOT NULL).
//   - Operation/Metadata: passthrough from caller; metadata is normalized to a non-nil map.
//   - Correlation: fills request_id and trace_id from context when available.
//
// OccurredAt is left empty; RecordChange will stamp the current time.
func BuildRecordFromContext(
	ctx context.Context,
	userID int64,
	target auditmodel.AuditTarget,
	operation auditmodel.AuditOperation,
	metadata map[string]any,
) auditmodel.RecordInput {
	actor := ActorFromContext(ctx, userID)

	if target.Version == nil {
		zero := int64(0)
		target.Version = &zero
	}

	if metadata == nil {
		metadata = make(map[string]any)
	}

	corr := auditmodel.AuditCorrelation{
		RequestID: utils.GetRequestIDFromContext(ctx),
	}
	if spanCtx := oteltrace.SpanFromContext(ctx).SpanContext(); spanCtx.IsValid() {
		corr.TraceID = spanCtx.TraceID().String()
	}

	return auditmodel.RecordInput{
		Actor:       actor,
		Target:      target,
		Operation:   operation,
		Metadata:    metadata,
		Correlation: corr,
	}
}
