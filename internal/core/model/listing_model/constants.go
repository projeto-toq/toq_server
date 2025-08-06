package listingmodel

type ListingStatus uint8

const (
	StatusDraft ListingStatus = iota + 1
	StatusAwaitingPhoto
	StatusAwaitingApproval
	StatusPublished
)

func (s ListingStatus) String() string {
	switch s {
	case StatusDraft:
		return "Draft"
	case StatusAwaitingPhoto:
		return "Awaiting Photo"
	case StatusAwaitingApproval:
		return "Awaiting Approval"
	case StatusPublished:
		return "Published"
	default:
		return "Unknown"
	}
}

type PropertyOwner uint8

const (
	OwnerMyself PropertyOwner = iota + 1
	OwnerSpouse
	OwnerParents
	OwnerGrandParents
	OwnerChildren
	OwnerUncles
	OwnerSiblings
)

func (o PropertyOwner) String() string {
	switch o {
	case OwnerMyself:
		return "Myself"
	case OwnerSpouse:
		return "Spouse"
	case OwnerParents:
		return "Parents"
	case OwnerGrandParents:
		return "GrandParents"
	case OwnerChildren:
		return "Children"
	case OwnerUncles:
		return "Uncles"
	case OwnerSiblings:
		return "Siblings"
	default:
		return "Unknown"
	}
}

type PropertyDelivered uint8

const (
	DeliveredFurnishedDecorated PropertyDelivered = iota + 1
	DeliveredFurnished
	DeliveredFixtured
	DeliverdAsPictured
)

func (d PropertyDelivered) String() string {
	switch d {
	case DeliveredFurnishedDecorated:
		return "Furnished & Decorated"
	case DeliveredFurnished:
		return "Furnished"
	case DeliveredFixtured:
		return "Fixtured"
	case DeliverdAsPictured:
		return "As Pictured"
	default:
		return "Unknown"
	}
}

type WhoLives uint8

const (
	LivesOwner WhoLives = iota + 1
	LivesTenant
	LivesVacant
)

func (w WhoLives) String() string {
	switch w {
	case LivesOwner:
		return "Owner"
	case LivesTenant:
		return "Tenant"
	case LivesVacant:
		return "Vacant"
	default:
		return "Unknown"
	}
}

type TransactionType uint8

const (
	TransactionSale TransactionType = iota + 1
	TransactionRent
	TransactionBoth
)

func (t TransactionType) String() string {
	switch t {
	case TransactionSale:
		return "Sale"
	case TransactionRent:
		return "Rent"
	case TransactionBoth:
		return "Both"
	default:
		return "Unknown"
	}
}

type InstallmentPlan uint8

const (
	PlanCash       InstallmentPlan = iota + 1 // a vista
	PlanShortTerm                             //0-6meses
	PlanMediumTerm                            //7-12meses
	PlanLongTerm                              //13-24meses
)

func (i InstallmentPlan) String() string {
	switch i {
	case PlanCash:
		return "Cash"
	case PlanShortTerm:
		return "Short Term"
	case PlanMediumTerm:
		return "Medium Term"
	case PlanLongTerm:
		return "Long Term"
	default:
		return "Unknown"
	}
}

type FinancingBlocker uint8

const (
	BlockerPendingProbate           FinancingBlocker = iota + 1 // Inventário (Probate)
	BlockerPendingLitigation                                    // Litígio (Litigation)
	BlockerPendingNoOccupancyPermit                             // Falta de Habite-se (No Occupancy Permit)
	Blockerother                                                // Outro (Other)
)

func (f FinancingBlocker) String() string {
	switch f {
	case BlockerPendingProbate:
		return "Pending Probate"
	case BlockerPendingLitigation:
		return "Pending Litigation"
	case BlockerPendingNoOccupancyPermit:
		return "Pending No Occupancy Permit"
	case Blockerother:
		return "Other"
	default:
		return "Unknown"
	}
}

type VisitType uint8

const (
	VisitClient VisitType = iota + 1
	VisitRealtor
	VisitContentProductor
	VisitAll
)

func (v VisitType) String() string {
	switch v {
	case VisitClient:
		return "Client"
	case VisitRealtor:
		return "Realtor"
	case VisitContentProductor:
		return "Content Productor"
	case VisitAll:
		return "All"
	default:
		return "Unknown"
	}
}

type AccompanyingType uint8

const (
	AccompanyingOwner AccompanyingType = iota + 1
	AccompanyingAssitant
	AccompanyingAlone
)

func (a AccompanyingType) String() string {
	switch a {
	case AccompanyingOwner:
		return "Owner"
	case AccompanyingAssitant:
		return "Assitant"
	case AccompanyingAlone:
		return "Alone"
	default:
		return "Unknown"
	}
}

type GuaranteeType uint8

const (
	GuaranteeDeposit    GuaranteeType = iota + 1 // Caução (Security Deposit)
	GuaranteeSurety                              // Fiança (Surety Bond/Guarantee)
	GuaranteeRentalBond                          // Seguro Fiança (Rental Insurance Bond)
)

func (g GuaranteeType) String() string {
	switch g {
	case GuaranteeDeposit:
		return "Security Deposit"
	case GuaranteeSurety:
		return "Surety Bond"
	case GuaranteeRentalBond:
		return "Rental Insurance Bond"
	default:
		return "Unknown"
	}
}
