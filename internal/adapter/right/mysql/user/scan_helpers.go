package mysqluseradapter

import (
	"database/sql"
	"fmt"

	permissionentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/entities"
	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
)

// scanUserWithRoleEntities scans multiple rows from a JOIN query (users + user_roles + roles)
// into strongly-typed UserWithRoleEntity structs.
//
// This function handles the complex scanning of 35 columns from the LEFT JOIN query,
// mapping each column to the appropriate entity field with proper NULL handling.
// It is specifically designed for queries that retrieve user data along with their
// active role information in a single database roundtrip.
//
// Used By:
//   - GetUserByID (internal/adapter/right/mysql/user/get_user_by_id.go)
//   - GetUserByNationalID (internal/adapter/right/mysql/user/get_user_by_nationalid.go)
//   - GetUserByPhoneNumber (internal/adapter/right/mysql/user/get_user_by_phone_number.go)
//
// Parameters:
//   - rows: SQL result set from GetUser* queries (caller must close)
//
// Returns:
//   - entities: Slice of UserWithRoleEntity (one per row, typically 0 or 1 row)
//   - error: Scanning errors (schema mismatch, type conversion failures)
//
// Column Order (MUST match query SELECT order exactly):
//
//	Columns 1-22: User fields (users table)
//	 1. id (INT UNSIGNED)
//	 2. full_name (VARCHAR)
//	 3. nick_name (VARCHAR, nullable)
//	 4. national_id (VARCHAR)
//	 5. creci_number (VARCHAR, nullable)
//	 6. creci_state (VARCHAR, nullable)
//	 7. creci_validity (DATE, nullable)
//	 8. born_at (DATE)
//	 9. phone_number (VARCHAR)
//	10. email (VARCHAR)
//	11. zip_code (VARCHAR)
//	12. street (VARCHAR)
//	13. number (VARCHAR)
//	14. complement (VARCHAR, nullable)
//	15. neighborhood (VARCHAR)
//	16. city (VARCHAR)
//	17. state (VARCHAR)
//	18. password (VARCHAR)
//	19. opt_status (TINYINT)
//	20. last_activity_at (TIMESTAMP)
//	21. deleted (TINYINT)
//
//	Columns 22-28: UserRole fields (user_roles table, ALL nullable due to LEFT JOIN)
//	22. ur.id (INT, nullable)
//	23. ur.user_id (INT, nullable)
//	24. ur.role_id (INT, nullable)
//	25. ur.is_active (TINYINT, nullable)
//	26. ur.status (TINYINT, nullable)
//	27. ur.expires_at (TIMESTAMP, nullable)
//	28. ur.blocked_until (DATETIME, nullable)
//
//	Columns 29-34: Role fields (roles table, ALL nullable due to LEFT JOIN)
//	29. r.id (INT, nullable)
//	30. r.slug (VARCHAR, nullable)
//	31. r.name (VARCHAR, nullable)
//	32. r.description (TEXT, nullable)
//	33. r.is_system_role (TINYINT, nullable)
//	34. r.is_active (TINYINT, nullable)
//
// NULL Handling:
//   - User fields: Uses sql.Null* types for optional fields (nick_name, creci_*, complement, etc.)
//   - UserRole fields: ALL nullable (user might not have active role)
//   - Role fields: ALL nullable (user might not have active role)
//   - LEFT JOIN ensures user is returned even without role (HasRole flag in entity)
//
// Performance:
//   - Single Scan() call per row (efficient memory usage)
//   - Type-safe scanning (no reflection overhead)
//   - Pre-allocates struct fields (no dynamic allocation during scan)
//
// Example Query That Uses This Scanner:
//
//	query := `SELECT
//	    u.id, u.full_name, u.nick_name, u.national_id, u.creci_number, u.creci_state, u.creci_validity,
//	    u.born_at, u.phone_number, u.email, u.zip_code, u.street, u.number, u.complement,
//	    u.neighborhood, u.city, u.state, u.password, u.opt_status, u.last_activity_at, u.deleted,
//	    ur.id, ur.user_id, ur.role_id, ur.is_active, ur.status, ur.expires_at, ur.blocked_until,
//	    r.id, r.slug, r.name, r.description, r.is_system_role, r.is_active
//	FROM users u
//	LEFT JOIN user_roles ur ON ur.user_id = u.id AND ur.is_active = 1
//	LEFT JOIN roles r ON r.id = ur.role_id
//	WHERE u.id = ? AND u.deleted = 0`
//
// Error Scenarios:
//   - Column count mismatch: Returns error with column index
//   - Type mismatch: Returns error with field name
//   - NULL in NOT NULL column: Returns scan error
func scanUserWithRoleEntities(rows *sql.Rows) ([]userentity.UserWithRoleEntity, error) {
	var entities []userentity.UserWithRoleEntity

	for rows.Next() {
		var entity userentity.UserWithRoleEntity

		// Scan all 34 columns from JOIN query
		// Order MUST match SELECT clause in get_user_by_*.go queries
		err := rows.Scan(
			// User fields (21 columns from users table)
			&entity.UserID,         // 1. u.id
			&entity.FullName,       // 2. u.full_name
			&entity.NickName,       // 3. u.nick_name (nullable)
			&entity.NationalID,     // 4. u.national_id
			&entity.CreciNumber,    // 5. u.creci_number (nullable)
			&entity.CreciState,     // 6. u.creci_state (nullable)
			&entity.CreciValidity,  // 7. u.creci_validity (nullable)
			&entity.BornAt,         // 8. u.born_at
			&entity.PhoneNumber,    // 9. u.phone_number
			&entity.Email,          // 10. u.email
			&entity.ZipCode,        // 11. u.zip_code
			&entity.Street,         // 12. u.street
			&entity.Number,         // 13. u.number
			&entity.Complement,     // 14. u.complement (nullable)
			&entity.Neighborhood,   // 15. u.neighborhood
			&entity.City,           // 16. u.city
			&entity.State,          // 17. u.state
			&entity.Password,       // 18. u.password
			&entity.OptStatus,      // 19. u.opt_status
			&entity.LastActivityAt, // 20. u.last_activity_at
			&entity.Deleted,        // 21. u.deleted

			// UserRole fields (7 columns from user_roles table, ALL nullable)
			&entity.UserRoleID,           // 22. ur.id
			&entity.UserRoleUserID,       // 23. ur.user_id
			&entity.UserRoleRoleID,       // 24. ur.role_id
			&entity.UserRoleIsActive,     // 25. ur.is_active
			&entity.UserRoleStatus,       // 26. ur.status
			&entity.UserRoleExpiresAt,    // 27. ur.expires_at
			&entity.UserRoleBlockedUntil, // 28. ur.blocked_until

			// Role fields (6 columns from roles table, ALL nullable)
			&entity.RoleID,           // 29. r.id
			&entity.RoleSlug,         // 30. r.slug
			&entity.RoleName,         // 31. r.name
			&entity.RoleDescription,  // 32. r.description
			&entity.RoleIsSystemRole, // 33. r.is_system_role
			&entity.RoleIsActive,     // 34. r.is_active
		)

		if err != nil {
			return nil, fmt.Errorf("scan user with role entity: %w", err)
		}

		entities = append(entities, entity)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return entities, nil
}

// scanUserRoleWithRoleEntities scans multiple rows from a JOIN query (user_roles + roles)
// into strongly-typed entities with embedded role data.
//
// This function handles scanning of 12 columns from the JOIN query, mapping each column
// to the appropriate entity field with proper NULL handling.
//
// Used By:
//   - GetUserRolesByUserID (internal/adapter/right/mysql/user/get_user_roles_by_user_id.go)
//
// Parameters:
//   - rows: SQL result set from GetUserRolesByUserID query (caller must close)
//
// Returns:
//   - userRoleEntities: Slice of UserRoleEntity (populated)
//   - roleEntities: Slice of RoleEntity (parallel array to userRoleEntities)
//   - error: Scanning errors (schema mismatch, type conversion failures)
//
// Column Order (MUST match query SELECT order exactly):
//
//	Columns 1-6: UserRole fields (user_roles table)
//	 1. ur.id (INT)
//	 2. ur.user_id (INT)
//	 3. ur.role_id (INT)
//	 4. ur.is_active (TINYINT)
//	 5. ur.status (TINYINT)
//	 6. ur.expires_at (TIMESTAMP, nullable)
//
//	Columns 7-12: Role fields (roles table)
//	 7. r.id (INT)
//	 8. r.slug (VARCHAR)
//	 9. r.name (VARCHAR)
//	10. r.description (TEXT, nullable)
//	11. r.is_system_role (TINYINT)
//	12. r.is_active (TINYINT)
//
// Performance:
//   - Single Scan() call per row (efficient memory usage)
//   - Type-safe scanning (no reflection overhead)
//   - Returns parallel arrays (index correspondence maintained)
//
// Example Query That Uses This Scanner:
//
//	query := `SELECT
//	    ur.id, ur.user_id, ur.role_id, ur.is_active, ur.status, ur.expires_at,
//	    r.id, r.slug, r.name, r.description, r.is_system_role, r.is_active
//	FROM user_roles ur
//	JOIN roles r ON r.id = ur.role_id
//	WHERE ur.user_id = ?`
//
// Error Scenarios:
//   - Column count mismatch: Returns error with column index
//   - Type mismatch: Returns error with field name
func scanUserRoleWithRoleEntities(rows *sql.Rows) ([]userentity.UserRoleEntity, []permissionentities.RoleEntity, error) {
	var userRoleEntities []userentity.UserRoleEntity
	var roleEntities []permissionentities.RoleEntity

	for rows.Next() {
		var (
			// UserRole fields
			id          int64
			userID      int64
			roleID      int64
			isActiveInt int64
			status      int64
			expiresAt   sql.NullTime

			// Role fields
			rID          int64
			rSlug        string
			rName        string
			rDescription sql.NullString
			rIsSystemInt int64
			rIsActiveInt int64
		)

		// Scan all 12 columns from JOIN query
		err := rows.Scan(
			// UserRole fields (6 columns)
			&id, &userID, &roleID, &isActiveInt, &status, &expiresAt,
			// Role fields (6 columns)
			&rID, &rSlug, &rName, &rDescription, &rIsSystemInt, &rIsActiveInt,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("scan user role with role entity: %w", err)
		}

		// Build UserRoleEntity
		userRoleEntity := userentity.UserRoleEntity{
			ID:       uint32(id),
			UserID:   uint32(userID),
			RoleID:   uint32(roleID),
			IsActive: isActiveInt == 1,
			Status:   int8(status),
		}
		if expiresAt.Valid {
			userRoleEntity.ExpiresAt = expiresAt
		}

		// Build RoleEntity
		roleEntity := permissionentities.RoleEntity{
			ID:   rID,
			Slug: rSlug,
			Name: rName,
			Description: func() string {
				if rDescription.Valid {
					return rDescription.String
				}
				return ""
			}(),
			IsSystemRole: rIsSystemInt == 1,
			IsActive:     rIsActiveInt == 1,
		}

		userRoleEntities = append(userRoleEntities, userRoleEntity)
		roleEntities = append(roleEntities, roleEntity)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return userRoleEntities, roleEntities, nil
}
