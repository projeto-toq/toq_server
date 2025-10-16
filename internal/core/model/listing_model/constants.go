package listingmodel

type ListingStatus uint8

const (
	// StatusDraft: O anúncio está sendo criado pelo proprietário e permanece invisível ao público.
	StatusDraft ListingStatus = iota + 1
	// StatusPendingPhotoScheduling: Anúncio criado e aguardando o agendamento da sessão de fotos.
	StatusPendingPhotoScheduling
	// StatusPhotosScheduled: Sessão de fotos agendada, aguardando execução.
	StatusPhotosScheduled
	// StatusPendingPhotoProcessing: Sessão concluída, aguardando tratamento e upload das fotos.
	StatusPendingPhotoProcessing
	// StatusPendingOwnerApproval: Materiais revisados e aguardando aprovação final do proprietário.
	StatusPendingOwnerApproval
	// StatusPendingAdminReview: Proprietário aprovou, aguardando revisão do time administrativo antes da publicação.
	StatusPendingAdminReview
	// StatusPublished: Anúncio ativo e visível publicamente.
	StatusPublished
	// StatusUnderOffer: Anúncio publicado que recebeu uma ou mais propostas.
	StatusUnderOffer
	// StatusUnderNegotiation: Uma proposta foi aceita e a negociação está em andamento.
	StatusUnderNegotiation
	// StatusClosed: O imóvel foi comercializado (vendido ou alugado) e o processo foi encerrado.
	StatusClosed
	// StatusSuspended: Anúncio pausado temporariamente pelo proprietário ou administrador.
	StatusSuspended
	// StatusRejectedByOwner: Versão final reprovada pelo proprietário (ex.: não aprovou as fotos).
	StatusRejectedByOwner
	// StatusNeedsRevision: Anúncio reprovado e aguardando ajustes antes de retornar ao fluxo de criação.
	StatusNeedsRevision
	// StatusExpired: Prazo de validade do anúncio encerrou sem negociação concluída.
	StatusExpired
	// StatusArchived: Anúncio removido do catálogo e mantido apenas para histórico.
	StatusArchived
)

func (s ListingStatus) String() string {
	switch s {
	case StatusDraft:
		return "DRAFT"
	case StatusPendingPhotoScheduling:
		return "PENDING_PHOTO_SCHEDULING"
	case StatusPhotosScheduled:
		return "PHOTOS_SCHEDULED"
	case StatusPendingPhotoProcessing:
		return "PENDING_PHOTO_PROCESSING"
	case StatusPendingOwnerApproval:
		return "PENDING_OWNER_APPROVAL"
	case StatusPendingAdminReview:
		return "PENDING_ADMIN_REVIEW"
	case StatusPublished:
		return "PUBLISHED"
	case StatusUnderOffer:
		return "UNDER_OFFER"
	case StatusUnderNegotiation:
		return "UNDER_NEGOTIATION"
	case StatusClosed:
		return "CLOSED"
	case StatusSuspended:
		return "SUSPENDED"
	case StatusRejectedByOwner:
		return "REJECTED_BY_OWNER"
	case StatusNeedsRevision:
		return "NEEDS_REVISION"
	case StatusExpired:
		return "EXPIRED"
	case StatusArchived:
		return "ARCHIVED"
	default:
		return "UNKNOWN"
	}
}

type PropertyOwner uint8
type PropertyDelivered uint8
type WhoLives uint8
type TransactionType uint8
type InstallmentPlan uint8
type FinancingBlocker uint8
type VisitType uint8
type AccompanyingType uint8
type GuaranteeType uint8

const (
	CatalogCategoryPropertyOwner     = "property_owner"
	CatalogCategoryPropertyDelivered = "property_delivered"
	CatalogCategoryWhoLives          = "who_lives"
	CatalogCategoryTransactionType   = "transaction_type"
	CatalogCategoryInstallmentPlan   = "installment_plan"
	CatalogCategoryFinancingBlocker  = "financing_blocker"
	CatalogCategoryVisitType         = "visit_type"
	CatalogCategoryAccompanyingType  = "accompanying_type"
	CatalogCategoryGuaranteeType     = "guarantee_type"
)
