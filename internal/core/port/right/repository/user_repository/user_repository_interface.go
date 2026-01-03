// Package userrepository exposes the UserRepoPortInterface contract mirrored by the MySQL adapter.
//
// Documentation rules (Section 8):
// - All methods are described with purpose, transactional expectations, returned errors, and edge cases.
// - Repositories return infrastructure errors (e.g., sql.ErrNoRows) and never map to HTTP.
// - Callers are responsible for starting/committing transactions when required.
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
	// ListUsersWithFilters returns a paginated admin listing with role info; read-only (tx optional).
	// Returns empty slice when no match; result.Total always set via COUNT DISTINCT.
	ListUsersWithFilters(ctx context.Context, tx *sql.Tx, filter ListUsersFilter) (ListUsersResult, error)
	// CreateAgencyInvite creates a new agency_invite; tx required for atomicity with related updates; returns constraint errors.
	CreateAgencyInvite(ctx context.Context, tx *sql.Tx, agency usermodel.UserInterface, phoneNumber string) (err error)
	// CreateAgencyRelationship links agency↔realtor in agency_realtor_relations; tx required; returns created relation id.
	CreateAgencyRelationship(ctx context.Context, tx *sql.Tx, agency usermodel.UserInterface, realtor usermodel.UserInterface) (id int64, err error)
	// CreateUser inserts a new user record; tx required; uniqueness violations bubble up; does not set roles.
	CreateUser(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (err error)
	// DeleteAgencyRealtorRelation removes agency↔realtor relation; tx required; returns rows deleted (sql.ErrNoRows if 0).
	DeleteAgencyRealtorRelation(ctx context.Context, tx *sql.Tx, agencyID int64, realtorID int64) (deleted int64, err error)
	// DeleteInviteByID removes agency_invite by id; tx required; returns rows deleted (sql.ErrNoRows if 0).
	DeleteInviteByID(ctx context.Context, tx *sql.Tx, id int64) (deleted int64, err error)
	// DeleteWrongSignInByUserID deletes temp_wrong_signin entry; tx required when part of unblock flow; returns rows deleted.
	DeleteWrongSignInByUserID(ctx context.Context, tx *sql.Tx, userID int64) (deleted int64, err error)
	// GetAgencyOfRealtor returns agency user for a realtor; tx optional; sql.ErrNoRows when relation not found.
	GetAgencyOfRealtor(ctx context.Context, tx *sql.Tx, realtorID int64) (agency usermodel.UserInterface, err error)
	// GetInviteByPhoneNumber finds active invite by phone; tx optional; sql.ErrNoRows when absent.
	GetInviteByPhoneNumber(ctx context.Context, tx *sql.Tx, phoneNumber string) (invite usermodel.InviteInterface, err error)
	// GetRealtorsByAgency lists realtor users linked to an agency; tx optional; returns empty slice when none.
	GetRealtorsByAgency(ctx context.Context, tx *sql.Tx, agencyID int64) (users []usermodel.UserInterface, err error)
	// GetUserByID fetches non-deleted user by id; tx optional; sql.ErrNoRows if not found.
	GetUserByID(ctx context.Context, tx *sql.Tx, id int64) (user usermodel.UserInterface, err error)
	// GetUserByNationalID fetches non-deleted user by CPF/CNPJ; tx optional; sql.ErrNoRows if not found.
	GetUserByNationalID(ctx context.Context, tx *sql.Tx, nationalID string) (user usermodel.UserInterface, err error)
	// GetUserByPhoneNumber fetches non-deleted user by phone; tx optional; sql.ErrNoRows if not found.
	GetUserByPhoneNumber(ctx context.Context, tx *sql.Tx, phoneNumber string) (user usermodel.UserInterface, err error)
	// ListAllUsers lists all non-deleted users; tx optional; sql.ErrNoRows when table has zero active rows.
	ListAllUsers(ctx context.Context, tx *sql.Tx) (users []usermodel.UserInterface, err error)
	// GetUsersByRoleAndStatus lists users with an active role matching slug/status; tx optional; empty slice when none.
	GetUsersByRoleAndStatus(ctx context.Context, tx *sql.Tx, role permissionmodel.RoleSlug, status globalmodel.UserRoleStatus) (users []usermodel.UserInterface, err error)
	// GetUserValidations returns temp_user_validations by user id; tx optional; sql.ErrNoRows when absent.
	GetUserValidations(ctx context.Context, tx *sql.Tx, id int64) (validation usermodel.ValidationInterface, err error)
	// GetWrongSigninByUserID returns temp_wrong_signin row; tx optional; sql.ErrNoRows when absent.
	GetWrongSigninByUserID(ctx context.Context, tx *sql.Tx, id int64) (wrongSignin usermodel.WrongSigninInterface, err error)
	// UpdateAgencyInviteByID updates phone/agency of an invite; tx required; sql.ErrNoRows if id not found.
	UpdateAgencyInviteByID(ctx context.Context, tx *sql.Tx, invite usermodel.InviteInterface) (err error)
	// UpdateUserByID updates user fields except password/last_activity; tx required; sql.ErrNoRows if not found/deleted.
	UpdateUserByID(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (err error)
	// UpdateUserLastActivity sets last_activity_at for a user; tx optional; sql.ErrNoRows if not found/deleted.
	UpdateUserLastActivity(ctx context.Context, tx *sql.Tx, id int64) (err error)
	// BatchUpdateUserLastActivity bulk-updates last_activity_at; tx optional; expects len(userIDs)==len(timestamps).
	BatchUpdateUserLastActivity(ctx context.Context, userIDs []int64, timestamps []int64) (err error)
	// UpdateUserPasswordByID updates password hash; tx required; sql.ErrNoRows if not found/deleted.
	UpdateUserPasswordByID(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (err error)
	// UpdateUserValidations upserts temp_user_validations; tx required; never returns sql.ErrNoRows.
	UpdateUserValidations(ctx context.Context, tx *sql.Tx, validation usermodel.ValidationInterface) (err error)
	// UpdateWrongSignIn upserts temp_wrong_signin (failed attempts); tx required; never returns sql.ErrNoRows.
	UpdateWrongSignIn(ctx context.Context, tx *sql.Tx, wrongSigin usermodel.WrongSigninInterface) (err error)
	// UpdateUserRoleStatusByUserID updates status of active role by user_id; tx optional; sql.ErrNoRows if no active role.
	UpdateUserRoleStatusByUserID(ctx context.Context, userID int64, status int) (err error)
	// UpdateUserRoleStatus updates status of active role filtered by slug; tx required; sql.ErrNoRows if no active role/slug.
	UpdateUserRoleStatus(ctx context.Context, tx *sql.Tx, userID int64, role permissionmodel.RoleSlug, status globalmodel.UserRoleStatus) error
	// ResetUserWrongSigninAttempts deletes temp_wrong_signin row; tx required alongside unblock; sql.ErrNoRows if already absent.
	ResetUserWrongSigninAttempts(ctx context.Context, tx *sql.Tx, userID int64) error
	// HasUserDuplicate checks uniqueness for cpf/email/phone (excluding current id); tx optional; returns true when duplicate exists.
	HasUserDuplicate(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) (exist bool, err error)
	// ExistsEmailForAnotherUser checks if email belongs to another user; tx optional; returns boolean without sql.ErrNoRows.
	ExistsEmailForAnotherUser(ctx context.Context, tx *sql.Tx, email string, excludeUserID int64) (bool, error)
	// ExistsPhoneForAnotherUser checks if phone belongs to another user; tx optional; returns boolean without sql.ErrNoRows.
	ExistsPhoneForAnotherUser(ctx context.Context, tx *sql.Tx, phone string, excludeUserID int64) (bool, error)
	// DeleteExpiredValidations deletes temp_user_validations past expiration with limit; tx required for controlled cleanup; returns rows deleted.
	DeleteExpiredValidations(ctx context.Context, tx *sql.Tx, limit int) (int64, error)

	// UserRole operations
	// CreateUserRole inserts user_role; tx required; returns created domain entity; constraint errors bubble up.
	CreateUserRole(ctx context.Context, tx *sql.Tx, userRole usermodel.UserRoleInterface) (usermodel.UserRoleInterface, error)
	// GetUserRolesByUserID lists all roles for a user; tx optional; empty slice when none.
	GetUserRolesByUserID(ctx context.Context, tx *sql.Tx, userID int64) ([]usermodel.UserRoleInterface, error)
	// GetActiveUserRoleByUserID fetches active role; tx optional; sql.ErrNoRows when no active role.
	GetActiveUserRoleByUserID(ctx context.Context, tx *sql.Tx, userID int64) (usermodel.UserRoleInterface, error)
	// GetUserRoleByUserIDAndRoleID fetches specific role mapping; tx optional; sql.ErrNoRows when absent.
	GetUserRoleByUserIDAndRoleID(ctx context.Context, tx *sql.Tx, userID, roleID int64) (usermodel.UserRoleInterface, error)
	// UpdateUserRole updates user_role fields; tx required; sql.ErrNoRows if id not found.
	UpdateUserRole(ctx context.Context, tx *sql.Tx, userRole usermodel.UserRoleInterface) error
	// DeleteUserRole hard-deletes user_role; tx required; sql.ErrNoRows if id not found.
	DeleteUserRole(ctx context.Context, tx *sql.Tx, userRoleID int64) error
	// DeactivateAllUserRoles sets is_active=0 for all roles of user; tx required; sql.ErrNoRows if none affected.
	DeactivateAllUserRoles(ctx context.Context, tx *sql.Tx, userID int64) error
	// ActivateUserRole sets is_active=1 for role/user pair; tx required; sql.ErrNoRows if pair not found.
	ActivateUserRole(ctx context.Context, tx *sql.Tx, userID, roleID int64) error

	// User blocking operations

	// SetUserBlockedUntil sets temporary block expiration (users.blocked_until); tx required; sql.ErrNoRows if user not found/deleted.
	SetUserBlockedUntil(ctx context.Context, tx *sql.Tx, userID int64, blockedUntil time.Time) error

	// ClearUserBlockedUntil sets blocked_until=NULL; tx required; sql.ErrNoRows if user not found/deleted.
	ClearUserBlockedUntil(ctx context.Context, tx *sql.Tx, userID int64) error

	// GetUsersWithExpiredBlock lists users with blocked_until <= NOW(); tx optional; empty slice when none.
	GetUsersWithExpiredBlock(ctx context.Context, tx *sql.Tx) ([]usermodel.UserInterface, error)

	// SetUserPermanentlyBlocked sets/clears permanently_blocked; tx required; clears blocked_until on unblock; sql.ErrNoRows if not found/deleted.
	SetUserPermanentlyBlocked(ctx context.Context, tx *sql.Tx, userID int64, blocked bool) error

	// ==================== Device Token Management ====================

	// AddDeviceToken upserts device_tokens scoped by user+device; tx optional; returns created/updated token; constraint errors bubble.
	AddDeviceToken(ctx context.Context, tx *sql.Tx, userID int64, deviceID, token string, platform *string) (usermodel.DeviceToken, error)

	// RemoveDeviceToken deletes device_tokens by user+token; tx optional; sql.ErrNoRows when token absent.
	RemoveDeviceToken(ctx context.Context, tx *sql.Tx, userID int64, token string) error

	// RemoveAllDeviceTokensByUserID bulk-deletes tokens for user; tx optional; 0 rows is success.
	RemoveAllDeviceTokensByUserID(ctx context.Context, tx *sql.Tx, userID int64) error

	// RemoveDeviceTokensByDeviceID bulk-deletes tokens for user+device; tx optional; 0 rows is success.
	RemoveDeviceTokensByDeviceID(ctx context.Context, tx *sql.Tx, userID int64, deviceID string) error

	// ListDeviceTokensByUserID returns device_tokens for user; tx optional; empty slice when none.
	ListDeviceTokensByUserID(ctx context.Context, tx *sql.Tx, userID int64) ([]usermodel.DeviceToken, error)

	// ListDeviceTokenStringsByUserIDIfOptedIn returns DISTINCT tokens for opted-in, non-deleted user; tx optional; empty slice otherwise.
	ListDeviceTokenStringsByUserIDIfOptedIn(ctx context.Context, tx *sql.Tx, userID int64) ([]string, error)

	// ListDeviceTokenStringsByOptedInUsers returns DISTINCT tokens for all opted-in, non-deleted users; tx optional; empty slice when none.
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
