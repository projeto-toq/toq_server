package usermodel

import "time"

const SystemUserID = 0

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
		"active",
		"blocked",
		"pending_profile",
		"pending_images",
		"pending_ocr",
		"reject_by_ocr",
		"pending_face",
		"reject_by_face",
		"pending_manual",
		"deleted",
		"invite_pending",
	}
	if us < StatusActive || int(us) >= len(statuses) {
		return "unknown"
	}
	return statuses[us]
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
