package userrepository

import (
	"context"
	"database/sql"
	"time"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

type UserRepoPortInterface interface {
	ListUsersWithFilters(ctx context.Context, tx *sql.Tx, filter ListUsersFilter) (ListUsersResult, error)
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
	ListAllUsers(ctx context.Context, tx *sql.Tx) (users []usermodel.UserInterface, err error)
	GetUsersByRoleAndStatus(ctx context.Context, tx *sql.Tx, role permissionmodel.RoleSlug, status globalmodel.UserRoleStatus) (users []usermodel.UserInterface, err error)
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
	UpdateUserRoleStatus(ctx context.Context, tx *sql.Tx, userID int64, role permissionmodel.RoleSlug, status globalmodel.UserRoleStatus) error
	ResetUserWrongSigninAttempts(ctx context.Context, tx *sql.Tx, userID int64) error
	HasUserDuplicate(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (exist bool, err error)
	ExistsEmailForAnotherUser(ctx context.Context, tx *sql.Tx, email string, excludeUserID int64) (bool, error)
	ExistsPhoneForAnotherUser(ctx context.Context, tx *sql.Tx, phone string, excludeUserID int64) (bool, error)
	DeleteExpiredValidations(ctx context.Context, tx *sql.Tx, limit int) (int64, error)

	// UserRole operations
	CreateUserRole(ctx context.Context, tx *sql.Tx, userRole usermodel.UserRoleInterface) (usermodel.UserRoleInterface, error)
	GetUserRolesByUserID(ctx context.Context, tx *sql.Tx, userID int64) ([]usermodel.UserRoleInterface, error)
	GetActiveUserRoleByUserID(ctx context.Context, tx *sql.Tx, userID int64) (usermodel.UserRoleInterface, error)
	GetUserRoleByUserIDAndRoleID(ctx context.Context, tx *sql.Tx, userID, roleID int64) (usermodel.UserRoleInterface, error)
	UpdateUserRole(ctx context.Context, tx *sql.Tx, userRole usermodel.UserRoleInterface) error
	DeleteUserRole(ctx context.Context, tx *sql.Tx, userRoleID int64) error
	DeactivateAllUserRoles(ctx context.Context, tx *sql.Tx, userID int64) error
	ActivateUserRole(ctx context.Context, tx *sql.Tx, userID, roleID int64) error

	// User blocking operations
	BlockUserTemporarily(ctx context.Context, tx *sql.Tx, userID int64, blockedUntil time.Time, reason string) error
	UnblockUser(ctx context.Context, tx *sql.Tx, userID int64) error
	GetExpiredTempBlockedUsers(ctx context.Context, tx *sql.Tx) ([]usermodel.UserRoleInterface, error)
}

type ListUsersFilter struct {
	Page             int
	Limit            int
	RoleName         string
	RoleSlug         string
	RoleStatus       *globalmodel.UserRoleStatus
	IsSystemRole     *bool
	FullName         string
	CPF              string
	Email            string
	PhoneNumber      string
	Deleted          *bool
	IDFrom           *int64
	IDTo             *int64
	BornAtFrom       *time.Time
	BornAtTo         *time.Time
	LastActivityFrom *time.Time
	LastActivityTo   *time.Time
}

type ListUsersResult struct {
	Users []usermodel.UserInterface
	Total int64
}
