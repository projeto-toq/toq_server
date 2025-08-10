package userrepository

import (
	"context"
	"database/sql"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

type UserRepoPortInterface interface {
	CreateAgencyInvite(ctx context.Context, tx *sql.Tx, agency usermodel.UserInterface, phoneNumber string) (err error)
	CreateAgencyRelationship(ctx context.Context, tx *sql.Tx, agency usermodel.UserInterface, realtor usermodel.UserInterface) (id int64, err error)
	CreateBaseRole(ctx context.Context, tx *sql.Tx, role usermodel.BaseRoleInterface) (err error)
	CreateUser(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (err error)
	CreateUserRole(ctx context.Context, tx *sql.Tx, role usermodel.UserRoleInterface) (err error)
	DeleteAgencyRealtorRelation(ctx context.Context, tx *sql.Tx, agencyID int64, realtorID int64) (deleted int64, err error)
	DeleteInviteByID(ctx context.Context, tx *sql.Tx, id int64) (deleted int64, err error)
	DeleteUserRolesByUserID(ctx context.Context, tx *sql.Tx, userID int64) (deleted int64, err error)
	DeleteWrongSignInByUserID(ctx context.Context, tx *sql.Tx, userID int64) (deleted int64, err error)
	GetAgencyOfRealtor(ctx context.Context, tx *sql.Tx, realtorID int64) (agency usermodel.UserInterface, err error)
	GetBaseRoleByRole(ctx context.Context, tx *sql.Tx, roleName usermodel.UserRole) (role usermodel.BaseRoleInterface, err error)
	GetBaseRoles(ctx context.Context, tx *sql.Tx) (roles []usermodel.BaseRoleInterface, err error)
	GetInviteByPhoneNumber(ctx context.Context, tx *sql.Tx, phoneNumber string) (invite usermodel.InviteInterface, err error)
	GetRealtorsByAgency(ctx context.Context, tx *sql.Tx, agencyID int64) (users []usermodel.UserInterface, err error)
	GetUserByID(ctx context.Context, tx *sql.Tx, id int64) (user usermodel.UserInterface, err error)
	GetUserByNationalID(ctx context.Context, tx *sql.Tx, nationalID string) (user usermodel.UserInterface, err error)
	GetUserByPhoneNumber(ctx context.Context, tx *sql.Tx, phoneNumber string) (user usermodel.UserInterface, err error)
	GetUserPhotoByID(ctx context.Context, tx *sql.Tx, id int64) (photo []byte, err error)
	GetUserRoleByRole(ctx context.Context, tx *sql.Tx, roleToGet usermodel.UserRole) (role usermodel.UserRoleInterface, err error)
	GetUserRolesByUserID(ctx context.Context, tx *sql.Tx, userID int64) (roles []usermodel.UserRoleInterface, err error)
	GetUsers(ctx context.Context, tx *sql.Tx) (users []usermodel.UserInterface, err error)
	GetUsersByStatus(ctx context.Context, tx *sql.Tx, UserRoleStatus usermodel.UserRoleStatus, userRole usermodel.UserRole) (users []usermodel.UserInterface, err error)
	GetUserValidations(ctx context.Context, tx *sql.Tx, id int64) (validation usermodel.ValidationInterface, err error)
	GetWrongSigninByUserID(ctx context.Context, tx *sql.Tx, id int64) (wrongSignin usermodel.WrongSigninInterface, err error)
	UpdateAgencyInviteByID(ctx context.Context, tx *sql.Tx, invite usermodel.InviteInterface) (err error)
	UpdateUserByID(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (err error)
	UpdateUserLastActivity(ctx context.Context, tx *sql.Tx, id int64) (err error)
	BatchUpdateUserLastActivity(ctx context.Context, userIDs []int64, timestamps []int64) (err error)
	UpdateUserPasswordByID(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (err error)
	UpdateUserPhotoByID(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (err error)
	UpdateUserRole(ctx context.Context, tx *sql.Tx, role usermodel.UserRoleInterface) (err error)
	UpdateUserValidations(ctx context.Context, tx *sql.Tx, validation usermodel.ValidationInterface) (err error)
	UpdateWrongSignIn(ctx context.Context, tx *sql.Tx, wrongSigin usermodel.WrongSigninInterface) (err error)
	VerifyUserDuplicity(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (exist bool, err error)
	AddDeviceToken(ctx context.Context, tx *sql.Tx, userID int64, token string, platform *string) error
	RemoveDeviceToken(ctx context.Context, tx *sql.Tx, userID int64, token string) error
	RemoveAllDeviceTokens(ctx context.Context, tx *sql.Tx, userID int64) error
}
