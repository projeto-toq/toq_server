package auditservice

import (
	"context"
	"time"

	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func (s *auditService) buildEvent(ctx context.Context, input auditmodel.RecordInput) auditmodel.AuditEvent {
	event := auditmodel.NewEvent()

	occurredAt := input.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now().UTC()
	}
	event.SetOccurredAt(occurredAt)
	event.SetActor(input.Actor)

	// Default version to zero to satisfy NOT NULL constraint when caller omits it.
	target := input.Target
	if target.Version == nil {
		zero := int64(0)
		target.Version = &zero
	}
	event.SetTarget(target)
	event.SetOperation(input.Operation)

	metadata := input.Metadata
	if metadata == nil {
		metadata = make(map[string]any)
	}
	event.SetMetadata(metadata)

	corr := input.Correlation
	if corr.RequestID == "" {
		corr.RequestID = utils.GetRequestIDFromContext(ctx)
	}
	if corr.TraceID == "" {
		spanCtx := oteltrace.SpanFromContext(ctx).SpanContext()
		if spanCtx.IsValid() {
			corr.TraceID = spanCtx.TraceID().String()
		}
	}
	event.SetCorrelation(corr)

	return event
}

// ActorFromContext builds an AuditActor using known context keys and an explicit user ID fallback.
func ActorFromContext(ctx context.Context, userID int64) auditmodel.AuditActor {
	actor := auditmodel.AuditActor{ID: userID}
	if ctx == nil {
		return actor
	}
	if infos, ok := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos); ok {
		actor.RoleSlug = string(infos.RoleSlug)
	}
	if v := ctx.Value(globalmodel.DeviceIDKey); v != nil {
		if deviceID, ok := v.(string); ok {
			actor.DeviceID = deviceID
		}
	}
	if v := ctx.Value(globalmodel.ClientIPKey); v != nil {
		if ip, ok := v.(string); ok {
			actor.IP = ip
		}
	}
	if v := ctx.Value(globalmodel.UserAgentKey); v != nil {
		if ua, ok := v.(string); ok {
			actor.UserAgent = ua
		}
	}
	return actor
}
