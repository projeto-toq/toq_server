package listingmodel

type ListingStatus uint8

const (
	// 1 StatusDraft: O anúncio está sendo criado pelo proprietário e permanece invisível ao público.
	StatusDraft ListingStatus = iota + 1
	// 2 StatusPendingAvailability: Anúncio criado, aguardando criação de agenda de disponibilidades do imóvel
	StatusPendingAvailability
	// 3 StatusPendingPhotoScheduling: Anúncio criado e aguardando o agendamento da sessão de fotos.
	StatusPendingPhotoScheduling
	// 4 StatusPendingPhotoConfirmation: Solicitado fotos para slot disponível. ag confirmação
	StatusPendingPhotoConfirmation
	// 5 StatusPhotosScheduled: Sessão de fotos agendada, aguardando execução.
	StatusPhotosScheduled
	// 6 StatusPendingPhotoProcessing: Sessão concluída, aguardando tratamento e upload das fotos.
	StatusPendingPhotoProcessing
	// 7 StatusPendingOwnerApproval: Materiais revisados e aguardando aprovação final do proprietário.
	StatusPendingOwnerApproval
	// 8 StatusRejectedByOwner: Versão final reprovada pelo proprietário (ex.: não aprovou as fotos).
	StatusRejectedByOwner
	// 9 StatusPendingAdminReview: Proprietário aprovou, aguardando revisão do time administrativo antes da publicação.
	StatusPendingAdminReview
	// 10 StatusReady: Anúncio aprovado e pronto para publicação (estado intermediário, não exposto externamente).
	StatusReady
	// 11 StatusPublished: Anúncio ativo e visível publicamente.
	StatusPublished
	// 12 StatusClosed: O imóvel foi comercializado (vendido ou alugado) e o processo foi encerrado.
	StatusClosed
	// 13 StatusSuspended: Anúncio pausado temporariamente pelo proprietário ou administrador.
	StatusSuspended
	// 14 StatusExpired: Prazo de validade do anúncio encerrou sem negociação concluída.
	StatusExpired
	// 15 StatusArchived: Anúncio removido do catálogo e mantido apenas para histórico.
	StatusArchived
	// 16 StatusNeedsRevision: Anúncio reprovado e aguardando ajustes antes de retornar ao fluxo de criação.
	StatusNeedsRevision
	// 17 StatusPendingPlanLoading: Anúncio de obra aguardando upload de plantas/renders de projeto.
	StatusPendingPlanLoading
)

func (s ListingStatus) String() string {
	if desc, ok := descriptorByStatus[s]; ok {
		return desc.Slug
	}
	return "UNKNOWN"
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
