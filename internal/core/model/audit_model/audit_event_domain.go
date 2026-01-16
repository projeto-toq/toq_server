package auditmodel

import "time"

// auditEvent is the concrete implementation of AuditEvent.
type auditEvent struct {
	id          int64
	occurredAt  time.Time
	actor       AuditActor
	target      AuditTarget
	operation   AuditOperation
	metadata    map[string]any
	correlation AuditCorrelation
}

// NewEvent creates a new audit event with zero values.
func NewEvent() AuditEvent {
	return &auditEvent{}
}

func (e *auditEvent) ID() int64 { return e.id }

func (e *auditEvent) SetID(id int64) { e.id = id }

func (e *auditEvent) OccurredAt() time.Time { return e.occurredAt }

func (e *auditEvent) SetOccurredAt(t time.Time) { e.occurredAt = t }

func (e *auditEvent) Actor() AuditActor { return e.actor }

func (e *auditEvent) SetActor(actor AuditActor) { e.actor = actor }

func (e *auditEvent) Target() AuditTarget { return e.target }

func (e *auditEvent) SetTarget(target AuditTarget) { e.target = target }

func (e *auditEvent) Operation() AuditOperation { return e.operation }

func (e *auditEvent) SetOperation(op AuditOperation) { e.operation = op }

func (e *auditEvent) Metadata() map[string]any { return e.metadata }

func (e *auditEvent) SetMetadata(m map[string]any) { e.metadata = m }

func (e *auditEvent) Correlation() AuditCorrelation { return e.correlation }

func (e *auditEvent) SetCorrelation(c AuditCorrelation) { e.correlation = c }
