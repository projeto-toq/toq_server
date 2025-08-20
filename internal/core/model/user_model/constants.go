package usermodel

import "time"

const SystemUserID = 0

type UserRole uint8

const (
	RoleRoot    UserRole = iota //0
	RoleOwner                   //1
	RoleRealtor                 //2
	RoleAgency                  //3
)

func (ur UserRole) String() string {
	roles := [...]string{
		"Root",
		"Owner",
		"Realtor",
		"Agency",
	}
	if ur < RoleRoot || int(ur) >= len(roles) {
		return "Unknown"
	}
	return roles[ur]
}

func UserRoleToUint8(userRole UserRole) uint8 {
	return uint8(userRole)
}

type GRPCService uint8

const (
	ServiceUserService GRPCService = iota
	ServiceListingService
	ServiceVisitService
	ServiceOfferService
)

func (gs GRPCService) String() string {
	services := [...]string{
		"UserService",
		"ListingService",
		"VisitService",
		"OfferService",
	}
	if gs < ServiceUserService || int(gs) >= len(services) {
		return "Unknown"
	}
	return services[gs]
}

func GrpcServiceToUint8(grpcService GRPCService) uint8 {
	return uint8(grpcService)
}

type UserRoleStatus int

const (
	StatusActive         UserRoleStatus = iota // normal user status
	StatusBlocked                              //nlocked by admin, see status reason
	StatusPendingProfile                       //awaiting phone and or email confirmation
	StatusPendingImages                        //awaiting creci images to be uploaded
	StatusPendingOCR                           //awaiting creci images to be OCR'd by AI
	StatusRejectByOCR                          //creci images were rejected by AI
	StatusPendingFace                          //awaiting face image to be verified by AI
	StatusRejectByFace                         //face image was rejected by AI
	StatusPendingManual                        //awaiting manual verification by admin
	StatusDeleted                              //user request the deletion of the account
	StatusInvitePending                        //realtor was invited and is pending acceptance
)

func (us UserRoleStatus) String() string {
	statuses := [...]string{
		"Active",
		"Blocked",
		"PendingProfile",
		"PendingImages",
		"PendingOCR",
		"RejectByOCR",
		"PendingFace",
		"RejectByFace",
		"PendingManual",
		"Deleted",
		"InvitePending",
	}
	if us < StatusActive || int(us) >= len(statuses) {
		return "Unknown"
	}
	return statuses[us]
}

type Privilege int

const (
	PrivilegeNone Privilege = iota
	PrivilegeCreate
	PrivilegeRead
	PrivilegeUpdate
	PrivilegeDelete
)

func (p Privilege) String() string {
	privileges := [...]string{
		"None",
		"Create",
		"Read",
		"Update",
		"Delete",
	}
	if p < PrivilegeNone || int(p) >= len(privileges) {
		return "Unknown"
	}
	return privileges[p]
}

const (
	// User validation codes
	EmailValidation = iota
	PhoneValidation
	PasswordValidation

	//expiraton time for validation codes
	ValidationCodeExpiration = 2 * time.Hour

	//user wrong signin attempts
	MaxWrongSigninAttempts = 3
)

type ActionFinished int

const (
	ActionFinishedCreated ActionFinished = iota
	ActionFinishedPhoneVerified
	ActionFinishedEmailVerified
	ActionFinishedCreciImagesUploaded
	ActionFinishedCreciNumberDoesntMatch
	ActionFinishedCreciStateDoesntMatch
	ActionFinishedCreciStateUnsupported
	ActionFinishedBadCreciImages
	ActionFinishedBadSelfieImage
	ActionFinishedSelfieDoesntMatch
	ActionFinishedCreciVerified
	ActionFinishedCreciFaceVerified
	ActionFinishedCreciManualVerified
	ActionFinishedCreciImagesUploadedForManualReview
	ActionFinishedInviteAccepted
	ActionFinishedInviteRejected
	ActionFinishedInviteCreated
)
