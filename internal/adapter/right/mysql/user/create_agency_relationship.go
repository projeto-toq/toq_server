package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) CreateAgencyRelationship(ctx context.Context, tx *sql.Tx, agency usermodel.UserInterface, realtor usermodel.UserInterface) (id int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	sql := `INSERT INTO realtors_agency (agency_id, realtor_id) VALUES (?, ?);`

	result, execErr := ua.ExecContext(ctx, tx, "insert", sql, agency.GetID(), realtor.GetID())
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.user.create_agency_relationship.exec_error", "error", execErr)
		return 0, fmt.Errorf("create agency relationship: %w", execErr)
	}

	id, lastErr := result.LastInsertId()
	if lastErr != nil {
		utils.SetSpanError(ctx, lastErr)
		logger.Error("mysql.user.create_agency_relationship.last_insert_id_error", "error", lastErr)
		return 0, fmt.Errorf("agency relationship last insert id: %w", lastErr)
	}

	return id, nil

}
