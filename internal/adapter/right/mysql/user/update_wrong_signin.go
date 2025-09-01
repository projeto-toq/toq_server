package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	userconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) UpdateWrongSignIn(ctx context.Context, tx *sql.Tx, wrongSigin usermodel.WrongSigninInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `INSERT INTO temp_wrong_signin (
				user_id, failed_attempts, last_attempt_at
				) VALUES (?, ?, ?)
				ON DUPLICATE KEY UPDATE
				failed_attempts = VALUES(failed_attempts),
				last_attempt_at = VALUES(last_attempt_at);`

	entity := userconverters.WrongSignInDomainToEntity(wrongSigin)

	_, err = ua.Update(ctx, tx, query,
		entity.UserID,
		entity.FailedAttempts,
		entity.LastAttemptAT,
	)
	if err != nil {
		slog.Error("mysqluseradapter/UpdateWrongSignIn: error executing Update", "error", err)
		return fmt.Errorf("update wrong_signin: %w", err)
	}

	return
}
