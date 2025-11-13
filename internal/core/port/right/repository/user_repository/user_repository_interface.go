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

	// SetUserBlockedUntil sets temporary block expiration for a user
	// Blocks user until specified timestamp (blocked_until column in users table)
	// Used by signin flow after MaxWrongSigninAttempts reached (configured in env.yaml)
	SetUserBlockedUntil(ctx context.Context, tx *sql.Tx, userID int64, blockedUntil time.Time) error

	// ClearUserBlockedUntil clears temporary block for a user
	// Sets blocked_until = NULL (unblocks user)
	// Used by worker when blocked_until expires, or by signin on success
	ClearUserBlockedUntil(ctx context.Context, tx *sql.Tx, userID int64) error

	// GetUsersWithExpiredBlock returns users whose blocked_until has passed
	// Used by worker to automatically unblock users
	// Returns empty slice if no expired blocks found (not an error)
	GetUsersWithExpiredBlock(ctx context.Context, tx *sql.Tx) ([]usermodel.UserInterface, error)

	// SetUserPermanentlyBlocked sets or clears permanent admin block
	// Used by admin endpoints to permanently block/unblock users
	SetUserPermanentlyBlocked(ctx context.Context, tx *sql.Tx, userID int64, blocked bool) error

	// ==================== Device Token Management ====================

	// AddDeviceToken registers a new push notification token for a user device
	// Uses INSERT ... ON DUPLICATE KEY UPDATE to handle token rotation for same device
	// Returns the created/updated device token record
	//
	// Parameters:
	//   - ctx: Context for tracing and cancellation
	//   - tx: Database transaction (can be nil for standalone operation)
	//   - userID: User's unique identifier
	//   - deviceID: Unique device identifier (UUIDv4)
	//   - token: FCM or APNs push notification token
	//   - platform: Device platform ("android"/"ios"/"web") - optional
	//
	// Returns:
	//   - token: Created device token record with ID
	//   - error: Database errors
	AddDeviceToken(ctx context.Context, tx *sql.Tx, userID int64, deviceID, token string, platform *string) (usermodel.DeviceToken, error)

	// RemoveDeviceToken deletes a specific push notification token
	// Returns sql.ErrNoRows if token not found
	//
	// Parameters:
	//   - ctx: Context for tracing
	//   - tx: Database transaction (can be nil)
	//   - userID: User's unique identifier
	//   - token: The token string to remove
	//
	// Returns:
	//   - error: sql.ErrNoRows if not found, or database errors
	RemoveDeviceToken(ctx context.Context, tx *sql.Tx, userID int64, token string) error

	// RemoveAllDeviceTokensByUserID removes all tokens for a user (cleanup on logout/delete)
	// Returns no error if user has no tokens (0 rows affected is success)
	//
	// Parameters:
	//   - ctx: Context for tracing
	//   - tx: Database transaction (can be nil)
	//   - userID: User's unique identifier
	//
	// Returns:
	//   - error: Database errors
	RemoveAllDeviceTokensByUserID(ctx context.Context, tx *sql.Tx, userID int64) error

	// RemoveDeviceTokensByDeviceID removes all tokens for a specific device
	// Used when device session is revoked or device is unregistered
	// Returns no error if device has no tokens
	//
	// Parameters:
	//   - ctx: Context for tracing
	//   - tx: Database transaction (can be nil)
	//   - userID: User's unique identifier
	//   - deviceID: Unique device identifier
	//
	// Returns:
	//   - error: Database errors
	RemoveDeviceTokensByDeviceID(ctx context.Context, tx *sql.Tx, userID int64, deviceID string) error

	// ListDeviceTokensByUserID retrieves all tokens for a user
	// Returns empty slice if user has no tokens (NOT sql.ErrNoRows)
	//
	// Parameters:
	//   - ctx: Context for tracing
	//   - tx: Database transaction (can be nil)
	//   - userID: User's unique identifier
	//
	// Returns:
	//   - tokens: Slice of device tokens
	//   - error: Database errors
	ListDeviceTokensByUserID(ctx context.Context, tx *sql.Tx, userID int64) ([]usermodel.DeviceToken, error)

	// ListDeviceTokenStringsByUserIDIfOptedIn retrieves FCM token strings for opted-in user
	// Only returns tokens if user.opt_status = 1 and user.deleted = 0
	// Returns empty slice if user opted out or deleted (NOT sql.ErrNoRows)
	//
	// Use Case: Targeted push notifications to single user
	//
	// Parameters:
	//   - ctx: Context for tracing
	//   - tx: Database transaction (can be nil)
	//   - userID: User's unique identifier
	//
	// Returns:
	//   - tokens: Slice of FCM token strings (DISTINCT)
	//   - error: Database errors
	ListDeviceTokenStringsByUserIDIfOptedIn(ctx context.Context, tx *sql.Tx, userID int64) ([]string, error)

	// ListDeviceTokenStringsByOptedInUsers retrieves all tokens for all opted-in users
	// Only returns tokens where user.opt_status = 1 and user.deleted = 0
	// Returns empty slice if no users opted in (NOT sql.ErrNoRows)
	//
	// Use Case: Broadcast notifications to all opted-in users
	// Warning: Can return large result set (thousands of tokens)
	//
	// Parameters:
	//   - ctx: Context for tracing
	//   - tx: Database transaction (can be nil)
	//
	// Returns:
	//   - tokens: Slice of FCM token strings (DISTINCT)
	//   - error: Database errors
	ListDeviceTokenStringsByOptedInUsers(ctx context.Context, tx *sql.Tx) ([]string, error)
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
