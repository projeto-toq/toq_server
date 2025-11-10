package userservices

import (
	"context"
	"database/sql"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetUserByID returns the domain user with the active role eagerly loaded.
//
// This method now delegates to the repository which performs an optimized JOIN query,
// returning the complete user aggregate (User + ActiveRole) in a single database round-trip.
//
// The service layer validates the domain invariant that every user MUST have an active role.
// If the repository returns a user without active role, this is treated as an error.
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - id: User's unique identifier
//
// Returns:
//   - user: UserInterface with ActiveRole populated
//   - error: Domain error with appropriate HTTP status code:
//   - 404 (Not Found) if user doesn't exist
//   - 500 (Internal) if active role is missing (domain invariant violation)
//   - 500 (Internal) for infrastructure failures (DB errors)
//
// Performance Improvement:
//   - Old: 2 queries (GetUserByID + GetActiveUserRoleByUserID)
//   - New: 1 query (repository JOIN)
//   - Latency reduction: ~50% (eliminates one DB round-trip)
func (us *userService) GetUserByID(ctx context.Context, id int64) (user usermodel.UserInterface, err error) {
	ctx, spanEnd, terr := utils.GenerateTracer(ctx)
	if terr != nil {
		return nil, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)

	// Iniciar transação somente leitura para leitura consistente
	tx, txErr := us.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		utils.LoggerFromContext(ctx).Error("user.get_by_id.tx_start_error", "error", txErr)
		return nil, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				utils.LoggerFromContext(ctx).Error("user.get_by_id.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	user, err = us.GetUserByIDWithTx(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if cmErr := us.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		utils.LoggerFromContext(ctx).Error("user.get_by_id.tx_commit_error", "error", cmErr)
		return nil, utils.InternalError("Failed to commit transaction")
	}

	return user, nil
}

// GetUserByIDWithTx loads the user with its active role using the provided transaction.
//
// IMPORTANT: After refactoring, this method is SIMPLIFIED. The repository now returns
// the complete user aggregate, so we no longer need to make a second query for active role.
//
// Domain Invariant Validation:
//   - Every user MUST have exactly one active role
//   - If repository returns user without active role, this is an error
//   - Service logs the violation and returns Internal Error
//
// Parameters:
//   - ctx: Context for logging
//   - tx: Database transaction (must not be nil)
//   - id: User's unique identifier
//
// Returns:
//   - user: UserInterface with ActiveRole populated
//   - error: Domain error (404 Not Found, 500 Internal)
func (us *userService) GetUserByIDWithTx(ctx context.Context, tx *sql.Tx, id int64) (user usermodel.UserInterface, err error) {
	ctx = utils.ContextWithLogger(ctx)

	// Repository now returns user WITH active role in single query
	user, err = us.repo.GetUserByID(ctx, tx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found or deleted
			return nil, utils.NotFoundError("User")
		}
		// Infrastructure error (DB failure, connection lost, etc.)
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.get_by_id.read_user_error", "error", err, "user_id", id)
		return nil, utils.InternalError("Failed to get user by ID")
	}

	// Validate domain invariant: every user MUST have active role
	if user.GetActiveRole() == nil {
		// This should NEVER happen if database is consistent
		// (every user should have at least one active role)
		utils.LoggerFromContext(ctx).Error("user.active_role.missing", "user_id", id)
		derr := utils.InternalError("User active role missing")
		utils.SetSpanError(ctx, derr)
		return nil, derr
	}

	return user, nil
}
