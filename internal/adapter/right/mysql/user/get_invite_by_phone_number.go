package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	userconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) GetInviteByPhoneNumber(ctx context.Context, tx *sql.Tx, phoneNumber string) (invite usermodel.InviteInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Query agency invitation by phone number
	// Note: phone_number is indexed but not unique, though business logic prevents duplicates per agency
	query := `SELECT id, agency_id, phone_number 
	          FROM agency_invites WHERE phone_number = ?;`

	rows, queryErr := ua.QueryContext(ctx, tx, "select", query, phoneNumber)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.user.get_invite_by_phone.query_error", "error", queryErr)
		return nil, fmt.Errorf("get invite by phone number query: %w", queryErr)
	}
	defer rows.Close()

	// Scan rows using type-safe function (replaces rowsToEntities)
	entities, err := scanInviteEntities(rows)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.get_invite_by_phone.scan_error", "error", err)
		return nil, fmt.Errorf("scan invite by phone rows: %w", err)
	}

	// Handle no results
	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	// Convert first entity to domain model using type-safe converter
	// Note: Query may return multiple rows if multiple agencies invited same number
	// Service layer should handle multi-agency invitation logic
	invite = userconverters.AgencyInviteEntityToDomainTyped(entities[0])

	return invite, nil
}
