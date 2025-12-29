package listingmodel

type ListingStatus uint8

const (
	// 0 StatusDraft: O anúncio está sendo criado pelo proprietário e permanece invisível ao público.
	StatusDraft ListingStatus = iota + 1
	// 1 StatusPendingAvailability: Anúncio criado, aguardando criação de agenda de disponibilidades do imóvel
	StatusPendingAvailability
	// 2 StatusPendingPhotoScheduling: Anúncio criado e aguardando o agendamento da sessão de fotos.
	StatusPendingPhotoScheduling
	// 3 StatusPendingPhotoConfirmation: Solicitado fotos para slot disponível. ag confirmação
	StatusPendingPhotoConfirmation
	// 4 StatusPhotosScheduled: Sessão de fotos agendada, aguardando execução.
	StatusPhotosScheduled
	// 5 StatusPendingPhotoProcessing: Sessão concluída, aguardando tratamento e upload das fotos.
	StatusPendingPhotoProcessing
	// 6 StatusPendingOwnerApproval: Materiais revisados e aguardando aprovação final do proprietário.
	StatusPendingOwnerApproval
	// 7 StatusRejectedByOwner: Versão final reprovada pelo proprietário (ex.: não aprovou as fotos).
	StatusRejectedByOwner
	// 8 StatusPendingAdminReview: Proprietário aprovou, aguardando revisão do time administrativo antes da publicação.
	StatusPendingAdminReview
	// 9 StatusReady: Anúncio aprovado e pronto para publicação (estado intermediário, não exposto externamente).
	StatusReady
	// 10 StatusPublished: Anúncio ativo e visível publicamente.
	StatusPublished
	// 11 StatusUnderOffer: Anúncio publicado que recebeu uma ou mais propostas.
	StatusUnderOffer
	// 12 StatusUnderNegotiation: Uma proposta foi aceita e a negociação está em andamento.
	StatusUnderNegotiation
	// 13 StatusClosed: O imóvel foi comercializado (vendido ou alugado) e o processo foi encerrado.
	StatusClosed
	// 14 StatusSuspended: Anúncio pausado temporariamente pelo proprietário ou administrador.
	StatusSuspended
	// 15 StatusExpired: Prazo de validade do anúncio encerrou sem negociação concluída.
	StatusExpired
	// 16 StatusArchived: Anúncio removido do catálogo e mantido apenas para histórico.
	StatusArchived
	// 17 StatusNeedsRevision: Anúncio reprovado e aguardando ajustes antes de retornar ao fluxo de criação.
	StatusNeedsRevision
)

func (s ListingStatus) String() string {
	switch s {
	case StatusDraft:
		return "DRAFT"
	case StatusPendingAvailability:
		return "PENDING_AVAILABILITY"
	case StatusPendingPhotoScheduling:
		return "PENDING_PHOTO_SCHEDULING"
	case StatusPendingPhotoConfirmation:
		return "PENDING_PHOTO_CONFIRMATION"
	case StatusPhotosScheduled:
		return "PHOTOS_SCHEDULED"
	case StatusPendingPhotoProcessing:
		return "PENDING_PHOTO_PROCESSING"
	case StatusPendingOwnerApproval:
		return "PENDING_OWNER_APPROVAL"
	case StatusRejectedByOwner:
		return "REJECTED_BY_OWNER"
	case StatusPendingAdminReview:
		return "PENDING_ADMIN_REVIEW"
	case StatusReady:
		return "READY"
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
	case StatusExpired:
		return "EXPIRED"
	case StatusArchived:
		return "ARCHIVED"
	case StatusNeedsRevision:
		return "NEEDS_REVISION"
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
type LandTerrainType uint8
type WarehouseSector uint8

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
	CatalogCategoryLandTerrainType   = "land_terrain_type"
	CatalogCategoryWarehouseSector   = "warehouse_sector"
)
