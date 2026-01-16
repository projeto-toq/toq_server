package auditconverters

import (
	"context"

	audientities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/audit/entities"
	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
)

// EventDomainToEntity converts a domain audit event into a persistence entity.
func EventDomainToEntity(_ context.Context, event auditmodel.AuditEvent) (audientities.AuditEventEntity, error) {
	entity := audientities.AuditEventEntity{
		ID:             event.ID(),
		OccurredAt:     event.OccurredAt(),
		ActorID:        event.Actor().ID,
		ActorRole:      event.Actor().RoleSlug,
		ActorDeviceID:  event.Actor().DeviceID,
		ActorIP:        event.Actor().IP,
		ActorUserAgent: event.Actor().UserAgent,
		TargetType:     string(event.Target().Type),
		TargetID:       event.Target().ID,
		TargetVersion:  0,
		Operation:      string(event.Operation()),
		Metadata:       event.Metadata(),
		RequestID:      event.Correlation().RequestID,
		TraceID:        event.Correlation().TraceID,
	}

	if event.Target().Version != nil {
		entity.TargetVersion = *event.Target().Version
	}

	return entity, nil
}
