package globalmodel

import "time"

const (

	//go routines
	ElapseTime = 60 * time.Second

	AppVersion = "2.0.0" //será atualizado no momento do build com o comando
// //go build -ldflags "-X main.Version=0.0.1" cmd/toq_api/toq_api.go

)

type TableName string

const (
	TableUsers         TableName = "users"
	TableAgencyInvites TableName = "agency_invites"
	TableRealtorAgency TableName = "realtor_agency"
	// TableBaseRoles mantido por compatibilidade de código legado. Schema atual usa 'roles'.
	TableBaseRoles TableName = "base_roles"
	TableRoles     TableName = "roles"
	TableUserRoles TableName = "user_roles"
	TableListings  TableName = "listings"
)

func (t TableName) String() string {
	return string(t)
}

// UserRoleStatus representa os possíveis status de um user_role
// NOTE: User blocking moved to users table (blocked_until, permanently_blocked)
// This enum now represents ONLY evolutionary role validation status
type UserRoleStatus int

const (
	StatusActive          UserRoleStatus = iota // role approved and operational 0
	StatusPendingBoth                           // awaiting both email and phone confirmation 1
	StatusPendingEmail                          // awaiting email confirmation 2
	StatusPendingPhone                          // awaiting phone confirmation 3
	StatusPendingCreci                          // awaiting creci images to be uploaded 4
	StatusPendingCnpj                           // awaiting cnpj images to be uploaded 5
	StatusPendingManual                         // awaiting manual verification by admin 6
	StatusRejected                              // admin rejected the documentation (legacy/general) 7
	StatusRefusedImage                          // refused due to image issues (e.g., unreadable/invalid) 8
	StatusRefusedDocument                       // refused due to document mismatch/invalidity 9
	StatusRefusedData                           // refused due to data inconsistency 10
)

// String implementa fmt.Stringer para UserRoleStatus
func (us UserRoleStatus) String() string {
	statuses := [...]string{
		"active",
		"pending_both",
		"pending_email",
		"pending_phone",
		"pending_creci",
		"pending_cnpj",
		"pending_manual",
		"rejected",
		"refused_image",
		"refused_document",
		"refused_data",
	}
	if us < StatusActive || int(us) >= len(statuses) {
		return "unknown"
	}
	return statuses[us]
}

// IsManualApprovalTarget verifies if the status is allowed for manual approval actions.
func IsManualApprovalTarget(status UserRoleStatus) bool {
	switch status {
	case StatusActive, StatusRejected, StatusRefusedImage, StatusRefusedDocument, StatusRefusedData:
		return true
	default:
		return false
	}
}

type NotificationType int

const (
	NotificationEmailChange NotificationType = iota
	NotificationPhoneChange
	NotificationPasswordChange
	NotificationCreciStateUnsupported
	NotificationInvalidCreciState
	NotificationInvalidCreciNumber
	NotificationBadSelfieImage
	NotificationBadCreciImages
	NotificationCreciValidated
	NotificationRealtorInviteSMS
	NotificationRealtorInvitePush
	NotificationInviteAccepted
	NotificationInviteRejected
	NotificationRealtorRemovedFromAgency
	NotificationAgencyRemovedFromRealtor
)

type PropertyType uint16

const (
	Apartment       PropertyType = 1 << iota //1 - apartamento
	CommercialStore                          //2 - loja
	CommercialFloor                          //4 - sala
	Suite                                    //8 - conjunto
	House                                    //16 - casa
	OffPlanHouse                             //32 - casa na planta
	ResidencialLand                          //64 - terreno residencial
	CommercialLand                           //128 - terreno comercial
	Building                                 //256 - prédio
	Warehouse                                //512 - galpão
)

func (p PropertyType) String() string {
	switch p {
	case Apartment:
		return "Apartment"
	case CommercialStore:
		return "Commercial Store"
	case CommercialFloor:
		return "Commercial Floor"
	case Suite:
		return "Suite"
	case House:
		return "House"
	case OffPlanHouse:
		return "Off Plan House"
	case ResidencialLand:
		return "Residencial Land"
	case CommercialLand:
		return "Commercial Land"
	case Building:
		return "Building"
	case Warehouse:
		return "Warehouse"
	default:
		return "Unknown"
	}
}
