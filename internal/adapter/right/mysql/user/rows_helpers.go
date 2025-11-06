package mysqluseradapter

import (
	"database/sql"
	"fmt"

	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
)

// // scanUserEntity scans a single *sql.Row into a UserEntity struct
// //
// // This function provides type-safe scanning of database rows, eliminating
// // the need for manual type assertions and reducing error-prone column ordering dependencies.
// //
// // Parameters:
// //   - row: *sql.Row from QueryRowContext
// //
// // Returns:
// //   - entity: UserEntity with all fields populated from database row
// //   - error: Scan errors or column mismatch errors
// //
// //nolint:unused // Reserved for future use with QueryRowContext methods
// func scanUserEntity(row *sql.Row) (userentity.UserEntity, error) {
// 	var entity userentity.UserEntity
// 	err := row.Scan(
// 		&entity.ID,
// 		&entity.FullName,
// 		&entity.NickName,
// 		&entity.NationalID,
// 		&entity.CreciNumber,
// 		&entity.CreciState,
// 		&entity.CreciValidity,
// 		&entity.BornAT,
// 		&entity.PhoneNumber,
// 		&entity.Email,
// 		&entity.ZipCode,
// 		&entity.Street,
// 		&entity.Number,
// 		&entity.Complement,
// 		&entity.Neighborhood,
// 		&entity.City,
// 		&entity.State,
// 		&entity.Password,
// 		&entity.OptStatus,
// 		&entity.LastActivityAT,
// 		&entity.Deleted,
// 		&entity.LastSignInAttempt,
// 	)
// 	return entity, err
// }

// scanUserEntities scans multiple rows into a slice of UserEntity structs
//
// Iterates through *sql.Rows and populates a typed slice, providing compile-time
// type safety and clear error handling.
//
// Parameters:
//   - rows: *sql.Rows from QueryContext (caller must close)
//
// Returns:
//   - entities: Slice of UserEntity with all rows scanned
//   - error: Scan errors, column mismatch, or rows.Err()
func scanUserEntities(rows *sql.Rows) ([]userentity.UserEntity, error) {
	var entities []userentity.UserEntity

	for rows.Next() {
		var entity userentity.UserEntity
		err := rows.Scan(
			&entity.ID,
			&entity.FullName,
			&entity.NickName,
			&entity.NationalID,
			&entity.CreciNumber,
			&entity.CreciState,
			&entity.CreciValidity,
			&entity.BornAT,
			&entity.PhoneNumber,
			&entity.Email,
			&entity.ZipCode,
			&entity.Street,
			&entity.Number,
			&entity.Complement,
			&entity.Neighborhood,
			&entity.City,
			&entity.State,
			&entity.Password,
			&entity.OptStatus,
			&entity.LastActivityAT,
			&entity.Deleted,
			&entity.LastSignInAttempt,
		)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}

// UserEntityWithRole combines user data with optional role information
// Used by ListUsersWithFilters to represent JOIN query results
type UserEntityWithRole struct {
	User             userentity.UserEntity
	HasRole          bool
	UserRoleID       int64
	RoleID           int64
	RoleStatus       int
	RoleIsActive     bool
	RoleName         string
	RoleSlug         string
	RoleDescription  string
	RoleIsSystemRole bool
	RoleActive       bool
}

// scanUserEntitiesWithRoles scans multiple rows containing user + role JOIN data
//
// This function handles the complex query result from ListUsersWithFilters which
// joins users, user_roles, and roles tables. It constructs both UserEntity and
// associated role data in a type-safe manner.
//
// Query Structure Expected:
//   - First 22 columns: user fields (see scanUserEntities)
//   - Columns 23-26: user_role fields (id, role_id, status, is_active)
//   - Columns 27-31: role fields (id, name, slug, description, is_system_role, is_active)
//
// Parameters:
//   - rows: *sql.Rows from JOIN query (caller must close)
//
// Returns:
//   - entities: Slice of UserEntityWithRole structs
//   - error: Scan errors, column mismatch, or rows.Err()
func scanUserEntitiesWithRoles(rows *sql.Rows) ([]UserEntityWithRole, error) {
	var entities []UserEntityWithRole

	for rows.Next() {
		var entity UserEntityWithRole
		var userRoleID sql.NullInt64
		var roleID sql.NullInt64
		var roleStatus sql.NullInt64
		var roleIsActive sql.NullInt64
		var roleEntityID sql.NullInt64
		var roleName sql.NullString
		var roleSlug sql.NullString
		var roleDescription sql.NullString
		var roleIsSystemRole sql.NullInt64
		var roleActive sql.NullInt64

		// Scan user fields (22 columns)
		err := rows.Scan(
			&entity.User.ID,
			&entity.User.FullName,
			&entity.User.NickName,
			&entity.User.NationalID,
			&entity.User.CreciNumber,
			&entity.User.CreciState,
			&entity.User.CreciValidity,
			&entity.User.BornAT,
			&entity.User.PhoneNumber,
			&entity.User.Email,
			&entity.User.ZipCode,
			&entity.User.Street,
			&entity.User.Number,
			&entity.User.Complement,
			&entity.User.Neighborhood,
			&entity.User.City,
			&entity.User.State,
			&entity.User.Password,
			&entity.User.OptStatus,
			&entity.User.LastActivityAT,
			&entity.User.Deleted,
			&entity.User.LastSignInAttempt,
			// Scan user_role fields (4 columns, may be NULL if no role)
			&userRoleID,
			&roleID,
			&roleStatus,
			&roleIsActive,
			// Scan role fields (6 columns, may be NULL if no role)
			&roleEntityID,
			&roleName,
			&roleSlug,
			&roleDescription,
			&roleIsSystemRole,
			&roleActive,
		)
		if err != nil {
			return nil, fmt.Errorf("scan user with role: %w", err)
		}

		// Populate UserRole if present (LEFT JOIN may return NULL)
		if userRoleID.Valid {
			entity.HasRole = true
			entity.UserRoleID = userRoleID.Int64
			entity.RoleID = roleID.Int64
			entity.RoleStatus = int(roleStatus.Int64)
			entity.RoleIsActive = roleIsActive.Int64 == 1

			// Populate Role details if present
			if roleEntityID.Valid {
				entity.RoleName = roleName.String
				entity.RoleSlug = roleSlug.String
				entity.RoleDescription = roleDescription.String
				entity.RoleIsSystemRole = roleIsSystemRole.Int64 == 1
				entity.RoleActive = roleActive.Int64 == 1
			}
		}

		entities = append(entities, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}

// scanValidationEntities scans multiple rows into a slice of UserValidationEntity structs
//
// Provides type-safe scanning for temp_user_validations queries, eliminating
// runtime panics from incorrect column indexing.
//
// Parameters:
//   - rows: *sql.Rows from QueryContext (caller must close)
//
// Returns:
//   - entities: Slice of UserValidationEntity with all rows scanned
//   - error: Scan errors, column mismatch, or rows.Err()
//
// Query Columns Expected (in order):
//  1. user_id (INT)
//  2. new_email (VARCHAR, nullable)
//  3. email_code (VARCHAR, nullable)
//  4. email_code_exp (TIMESTAMP, nullable)
//  5. new_phone (VARCHAR, nullable)
//  6. phone_code (VARCHAR, nullable)
//  7. phone_code_exp (TIMESTAMP, nullable)
//  8. password_code (VARCHAR, nullable)
//  9. password_code_exp (TIMESTAMP, nullable)
func scanValidationEntities(rows *sql.Rows) ([]userentity.UserValidationEntity, error) {
	var entities []userentity.UserValidationEntity

	for rows.Next() {
		var entity userentity.UserValidationEntity
		err := rows.Scan(
			&entity.UserID,
			&entity.NewEmail,
			&entity.EmailCode,
			&entity.EmailCodeExp,
			&entity.NewPhone,
			&entity.PhoneCode,
			&entity.PhoneCodeExp,
			&entity.PasswordCode,
			&entity.PasswordCodeExp,
		)
		if err != nil {
			return nil, fmt.Errorf("scan validation entity: %w", err)
		}
		entities = append(entities, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("validation rows iteration error: %w", err)
	}

	return entities, nil
}

// scanInviteEntities scans multiple rows into a slice of AgencyInvite structs
//
// Provides type-safe scanning for agency_invites queries.
//
// Parameters:
//   - rows: *sql.Rows from QueryContext (caller must close)
//
// Returns:
//   - entities: Slice of AgencyInvite with all rows scanned
//   - error: Scan errors, column mismatch, or rows.Err()
//
// Query Columns Expected (in order):
//  1. id (INT, PRIMARY KEY)
//  2. agency_id (INT, FOREIGN KEY to users.id)
//  3. phone_number (VARCHAR)
func scanInviteEntities(rows *sql.Rows) ([]userentity.AgencyInvite, error) {
	var entities []userentity.AgencyInvite

	for rows.Next() {
		var entity userentity.AgencyInvite
		err := rows.Scan(
			&entity.ID,
			&entity.AgencyID,
			&entity.PhoneNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("scan invite entity: %w", err)
		}
		entities = append(entities, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("invite rows iteration error: %w", err)
	}

	return entities, nil
}

// scanWrongSigninEntities scans multiple rows into a slice of WrongSignInEntity structs
//
// Provides type-safe scanning for temp_wrong_signin queries.
//
// Parameters:
//   - rows: *sql.Rows from QueryContext (caller must close)
//
// Returns:
//   - entities: Slice of WrongSignInEntity with all rows scanned
//   - error: Scan errors, column mismatch, or rows.Err()
//
// Query Columns Expected (in order):
//  1. user_id (INT, PRIMARY KEY, FOREIGN KEY to users.id)
//  2. failed_attempts (TINYINT)
//  3. last_attempt_at (TIMESTAMP)
func scanWrongSigninEntities(rows *sql.Rows) ([]userentity.WrongSignInEntity, error) {
	var entities []userentity.WrongSignInEntity

	for rows.Next() {
		var entity userentity.WrongSignInEntity
		err := rows.Scan(
			&entity.UserID,
			&entity.FailedAttempts,
			&entity.LastAttemptAT,
		)
		if err != nil {
			return nil, fmt.Errorf("scan wrong signin entity: %w", err)
		}
		entities = append(entities, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("wrong signin rows iteration error: %w", err)
	}

	return entities, nil
}
