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
	TableBaseRoles     TableName = "base_roles"
	TableUserRoles     TableName = "user_roles"
	TableListings      TableName = "listings"
)

func (t TableName) String() string {
	return string(t)
}

// type AuditAction int

// const (
// 	AuditCreate AuditAction = iota
// 	AuditRead
// 	AuditUpdate
// 	AuditDelete
// )

// func (a AuditAction) String() string {
// 	switch a {
// 	case AuditCreate:
// 		return "Create"
// 	case AuditRead:
// 		return "Read"
// 	case AuditUpdate:
// 		return "Update"
// 	case AuditDelete:
// 		return "Delete"
// 	default:
// 		return "Unknown"
// 	}
// }

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
