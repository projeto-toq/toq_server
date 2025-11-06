package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateAgencyRelationship establishes a many-to-many relationship between agency and realtor
//
// This function creates a record in the realtors_agency junction table, linking
// a realtor user to an agency user. This relationship enables realtors to work
// under an agency's umbrella for commission splitting and listings management.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (REQUIRED - relationship must be transactional)
//   - agency: UserInterface representing the agency
//   - realtor: UserInterface representing the realtor being added
//
// Returns:
//   - id: Auto-generated primary key of the created relationship
//   - error: Database errors (duplicate relationship, foreign key violations)
//
// Business Rules:
//   - Both agency and realtor must exist (foreign key constraints)
//   - Duplicate relationships are allowed (no unique constraint on agency_id+realtor_id)
//   - Service layer validates that agency role is "agency" and realtor role is "realtor"
//
// Database Constraints:
//   - FK: agency_id REFERENCES users(id) ON DELETE CASCADE
//   - FK: realtor_id REFERENCES users(id) (no cascade)
//
// Usage Example:
//
//	relationID, err := adapter.CreateAgencyRelationship(ctx, tx, agency, realtor)
//	if err != nil {
//	    // Handle duplicate or FK violation
//	}
func (ua *UserAdapter) CreateAgencyRelationship(ctx context.Context, tx *sql.Tx, agency usermodel.UserInterface, realtor usermodel.UserInterface) (id int64, err error) {
	// Initialize tracing for observability
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Attach logger to context for request_id/trace_id propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Insert new agency-realtor relationship
	sql := `INSERT INTO realtors_agency (agency_id, realtor_id) VALUES (?, ?);`

	// Execute insert using instrumented adapter
	result, execErr := ua.ExecContext(ctx, tx, "insert", sql, agency.GetID(), realtor.GetID())
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.create_agency_relationship.exec_error", "error", execErr)
		return 0, fmt.Errorf("create agency relationship: %w", execErr)
	}

	// Retrieve auto-generated primary key
	id, lastErr := result.LastInsertId()
	if lastErr != nil {
		utils.SetSpanError(ctx, lastErr)
		logger.Error("mysql.user.create_agency_relationship.last_insert_id_error", "error", lastErr)
		return 0, fmt.Errorf("agency relationship last insert id: %w", lastErr)
	}

	return id, nil

}
