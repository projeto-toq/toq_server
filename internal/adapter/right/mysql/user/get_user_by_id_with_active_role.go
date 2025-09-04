package mysqluseradapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetUserByIDWithActiveRole loads the user and populates the active role (slug and status) if present.
// It runs within the provided transaction and returns a domain user with ActiveRole set.
func (ua *UserAdapter) GetUserByIDWithActiveRole(ctx context.Context, tx *sql.Tx, id int64) (usermodel.UserInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	// 1) Load base user
	user, err := ua.GetUserByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	// 2) Load active role (slug + status) using typed scan for reliability
	// Nota: ajustar para usar a tabela 'roles' (schema atual) em vez de 'base_roles'.
	query := `
		SELECT r.slug, ur.status
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_id = ? AND ur.is_active = 1
		LIMIT 1`

	var slug string
	var statusInt int64
	row := tx.QueryRowContext(ctx, query, id)
	if err := row.Scan(&slug, &statusInt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// no active role — return user without active role set; caller decides how to handle
			return user, nil
		}
		slog.Error("mysqluseradapter/GetUserByIDWithActiveRole: error scanning row", "error", err)
		return nil, fmt.Errorf("get active role scan: %w", err)
	}

	// Build and attach active role to user domain (permission model types)
	role := permissionmodel.NewRole()
	role.SetSlug(slug)

	userRole := permissionmodel.NewUserRole()
	userRole.SetRole(role)
	userRole.SetStatus(permissionmodel.UserRoleStatus(statusInt))

	user.SetActiveRole(userRole)
	return user, nil
}
