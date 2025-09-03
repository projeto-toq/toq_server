package userrepository

import (
	"context"
	"database/sql"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

type UserRepoPortInterface interface {
	CreateAgencyInvite(ctx context.Context, tx *sql.Tx, agency usermodel.UserInterface, phoneNumber string) (err error)
	CreateAgencyRelationship(ctx context.Context, tx *sql.Tx, agency usermodel.UserInterface, realtor usermodel.UserInterface) (id int64, err error)
	CreateUser(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (err error)
	DeleteAgencyRealtorRelation(ctx context.Context, tx *sql.Tx, agencyID int64, realtorID int64) (deleted int64, err error)
	DeleteInviteByID(ctx context.Context, tx *sql.Tx, id int64) (deleted int64, err error)
	DeleteWrongSignInByUserID(ctx context.Context, tx *sql.Tx, userID int64) (deleted int64, err error)
	GetAgencyOfRealtor(ctx context.Context, tx *sql.Tx, realtorID int64) (agency usermodel.UserInterface, err error)
	GetInviteByPhoneNumber(ctx context.Context, tx *sql.Tx, phoneNumber string) (invite usermodel.InviteInterface, err error)
	GetRealtorsByAgency(ctx context.Context, tx *sql.Tx, agencyID int64) (users []usermodel.UserInterface, err error)
	GetUserByID(ctx context.Context, tx *sql.Tx, id int64) (user usermodel.UserInterface, err error)
	GetUserByNationalID(ctx context.Context, tx *sql.Tx, nationalID string) (user usermodel.UserInterface, err error)
	GetUserByPhoneNumber(ctx context.Context, tx *sql.Tx, phoneNumber string) (user usermodel.UserInterface, err error)
	GetUsers(ctx context.Context, tx *sql.Tx) (users []usermodel.UserInterface, err error)
	GetUserValidations(ctx context.Context, tx *sql.Tx, id int64) (validation usermodel.ValidationInterface, err error)
	GetWrongSigninByUserID(ctx context.Context, tx *sql.Tx, id int64) (wrongSignin usermodel.WrongSigninInterface, err error)
	UpdateAgencyInviteByID(ctx context.Context, tx *sql.Tx, invite usermodel.InviteInterface) (err error)
	UpdateUserByID(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (err error)
	UpdateUserLastActivity(ctx context.Context, tx *sql.Tx, id int64) (err error)
	BatchUpdateUserLastActivity(ctx context.Context, userIDs []int64, timestamps []int64) (err error)
	UpdateUserPasswordByID(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (err error)
	UpdateUserValidations(ctx context.Context, tx *sql.Tx, validation usermodel.ValidationInterface) (err error)
	UpdateWrongSignIn(ctx context.Context, tx *sql.Tx, wrongSigin usermodel.WrongSigninInterface) (err error)
	UpdateUserRoleStatusByUserID(ctx context.Context, userID int64, status int) (err error)
	ResetUserWrongSigninAttempts(ctx context.Context, userID int64) (err error)
	VerifyUserDuplicity(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (exist bool, err error)
	// ExistsEmailForAnotherUser checks if an email is already used by a different user (deleted=0)
	ExistsEmailForAnotherUser(ctx context.Context, tx *sql.Tx, email string, excludeUserID int64) (bool, error)
	// ExistsPhoneForAnotherUser checks if a phone number is already used by a different user (deleted=0)
	ExistsPhoneForAnotherUser(ctx context.Context, tx *sql.Tx, phone string, excludeUserID int64) (bool, error)
	AddDeviceToken(ctx context.Context, tx *sql.Tx, userID int64, token string, platform *string) error
	RemoveDeviceToken(ctx context.Context, tx *sql.Tx, userID int64, token string) error
	RemoveAllDeviceTokens(ctx context.Context, tx *sql.Tx, userID int64) error

	// Per-device operations (backward-compatible when schema lacks device_id)
	AddTokenForDevice(ctx context.Context, tx *sql.Tx, userID int64, deviceID, token string, platform *string) error
	RemoveTokensByDeviceID(ctx context.Context, tx *sql.Tx, userID int64, deviceID string) error
}
