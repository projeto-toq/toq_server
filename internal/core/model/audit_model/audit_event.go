package auditmodel

import "time"

// AuditOperation represents the type of action being recorded in the audit trail.
type AuditOperation string

const (
	OperationCreate         AuditOperation = "create"
	OperationUpdate         AuditOperation = "update"
	OperationPromote        AuditOperation = "promote"
	OperationStatusChange   AuditOperation = "status_change"
	OperationPublish        AuditOperation = "publish"
	OperationUnpublish      AuditOperation = "unpublish"
	OperationDiscard        AuditOperation = "discard"
	OperationDelete         AuditOperation = "delete"
	OperationProposalCreate AuditOperation = "proposal_create"
	OperationProposalAccept AuditOperation = "proposal_accept"
	OperationProposalReject AuditOperation = "proposal_reject"
	OperationProposalCancel AuditOperation = "proposal_cancel"
	OperationVisitRequest   AuditOperation = "visit_request"
	OperationVisitApprove   AuditOperation = "visit_approve"
	OperationVisitReject    AuditOperation = "visit_reject"
	OperationVisitCancel    AuditOperation = "visit_cancel"
	OperationVisitComplete  AuditOperation = "visit_complete"
	OperationVisitNoShow    AuditOperation = "visit_no_show"
	OperationMediaApprove   AuditOperation = "media_approve"
	OperationMediaReject    AuditOperation = "media_reject"
	OperationAgendaCreate   AuditOperation = "agenda_create"
	OperationAgendaFinish   AuditOperation = "agenda_finish"
	OperationAuthSignin     AuditOperation = "auth_signin"
	OperationAuthSignout    AuditOperation = "auth_signout"
	OperationPasswordReset  AuditOperation = "password_reset"
)

// TargetType represents the audited resource domain.
type TargetType string

const (
	TargetListingIdentity TargetType = "listing_identities"
	TargetListingVersion  TargetType = "listing_versions"
	TargetListingVisit    TargetType = "listing_visits"
	TargetProposal        TargetType = "proposals"
	TargetMediaAsset      TargetType = "media_assets"
	TargetListingAgenda   TargetType = "listing_agendas"
	TargetSession         TargetType = "sessions"
)

// AuditActor identifies who performed the action.
type AuditActor struct {
	ID        int64
	RoleSlug  string
	DeviceID  string
	IP        string
	UserAgent string
}

// AuditTarget identifies the resource that changed.
type AuditTarget struct {
	Type    TargetType
	ID      int64
	Version *int64
}

// AuditCorrelation carries tracing/log correlation fields.
type AuditCorrelation struct {
	RequestID string
	TraceID   string
}

// AuditEvent defines the immutable audit trail entry.
type AuditEvent interface {
	ID() int64
	SetID(int64)
	OccurredAt() time.Time
	SetOccurredAt(time.Time)
	Actor() AuditActor
	SetActor(AuditActor)
	Target() AuditTarget
	SetTarget(AuditTarget)
	Operation() AuditOperation
	SetOperation(AuditOperation)
	Metadata() map[string]any
	SetMetadata(map[string]any)
	Correlation() AuditCorrelation
	SetCorrelation(AuditCorrelation)
}

// RecordInput is the DTO accepted by the audit service to register a new event.
type RecordInput struct {
	Actor       AuditActor
	Target      AuditTarget
	Operation   AuditOperation
	Metadata    map[string]any
	OccurredAt  time.Time
	Correlation AuditCorrelation
}
