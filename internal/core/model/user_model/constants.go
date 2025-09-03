package usermodel

import "time"

const SystemUserID = 0

const (
	// User validation codes
	EmailValidation = iota
	PhoneValidation
	PasswordValidation

	//expiraton time for validation codes
	ValidationCodeExpiration = 2 * time.Hour

	//user wrong signin attempts
	MaxWrongSigninAttempts = 3
	//temporary block duration for failed signin attempts
	TempBlockDuration = 15 * time.Minute
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
	// Profile verification aggregated actions (derived in services)
	ActionProfileEmailVerifiedPhonePending
	ActionProfilePhoneVerifiedEmailPending
	ActionProfileVerificationCompleted
)
