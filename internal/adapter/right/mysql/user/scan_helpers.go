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
//
// Parameters:
//   - rows: SQL result set from GetUserByID/GetUserByNationalID/GetUserByPhoneNumber queries
//
// Returns:
//   - entities: Slice of UserWithRoleEntity (one per row, typically 0 or 1 row)
//   - error: Scanning errors (schema mismatch, type conversion failures)
//
// Column Order (must match query SELECT order):
//   - Columns 1-22: User fields (users table)
//   - Columns 23-29: UserRole fields (user_roles table, nullable)
//   - Columns 30-35: Role fields (roles table, nullable)
//
// NULL Handling:
//   - User fields: Uses sql.Null* types for optional fields (nick_name, creci_*, complement, etc.)
//   - UserRole fields: ALL nullable (user might not have active role)
//   - Role fields: ALL nullable (user might not have active role)
//
// Performance:
//   - Single Scan() call per row (efficient memory usage)
//   - Type-safe scanning (no reflection overhead)
func scanUserWithRoleEntities(rows *sql.Rows) ([]userentity.UserWithRoleEntity, error) {
	var entities []userentity.UserWithRoleEntity

	for rows.Next() {
		var entity userentity.UserWithRoleEntity

		// Scan all 35 columns from JOIN query
		// Order MUST match SELECT clause in get_user_by_*.go queries
		err := rows.Scan(
			// User fields (22 columns from users table)
			&entity.UserID,            // 1. u.id
			&entity.FullName,          // 2. u.full_name
			&entity.NickName,          // 3. u.nick_name (nullable)
			&entity.NationalID,        // 4. u.national_id
			&entity.CreciNumber,       // 5. u.creci_number (nullable)
			&entity.CreciState,        // 6. u.creci_state (nullable)
			&entity.CreciValidity,     // 7. u.creci_validity (nullable)
			&entity.BornAt,            // 8. u.born_at
			&entity.PhoneNumber,       // 9. u.phone_number
			&entity.Email,             // 10. u.email
			&entity.ZipCode,           // 11. u.zip_code
			&entity.Street,            // 12. u.street
			&entity.Number,            // 13. u.number
			&entity.Complement,        // 14. u.complement (nullable)
			&entity.Neighborhood,      // 15. u.neighborhood
			&entity.City,              // 16. u.city
			&entity.State,             // 17. u.state
			&entity.Password,          // 18. u.password
			&entity.OptStatus,         // 19. u.opt_status
			&entity.LastActivityAt,    // 20. u.last_activity_at
			&entity.Deleted,           // 21. u.deleted
			&entity.LastSignInAttempt, // 22. u.last_signin_attempt (nullable)

			// UserRole fields (7 columns from user_roles table, ALL nullable)
			&entity.UserRoleID,           // 23. ur.id
			&entity.UserRoleUserID,       // 24. ur.user_id
			&entity.UserRoleRoleID,       // 25. ur.role_id
			&entity.UserRoleIsActive,     // 26. ur.is_active
			&entity.UserRoleStatus,       // 27. ur.status
			&entity.UserRoleExpiresAt,    // 28. ur.expires_at
			&entity.UserRoleBlockedUntil, // 29. ur.blocked_until

			// Role fields (6 columns from roles table, ALL nullable)
			&entity.RoleID,           // 30. r.id
			&entity.RoleSlug,         // 31. r.slug
			&entity.RoleName,         // 32. r.name
			&entity.RoleDescription,  // 33. r.description
			&entity.RoleIsSystemRole, // 34. r.is_system_role
			&entity.RoleIsActive,     // 35. r.is_active
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
// Parameters:
//   - rows: SQL result set from GetUserRolesByUserID query
//
// Returns:
//   - userRoleEntities: Slice of UserRoleEntity (populated)
//   - roleEntities: Slice of RoleEntity (parallel array to userRoleEntities)
//   - error: Scanning errors (schema mismatch, type conversion failures)
//
// Column Order (must match query SELECT order):
//   - Columns 1-6: UserRole fields (user_roles table)
//   - Columns 7-12: Role fields (roles table)
//
// Performance:
//   - Single Scan() call per row (efficient memory usage)
//   - Type-safe scanning (no reflection overhead)
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
