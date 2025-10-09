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

	id, err = ua.Create(ctx, tx, sql, agency.GetID(), realtor.GetID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.create_agency_relationship.create_error", "error", err)
		return 0, fmt.Errorf("create agency relationship: %w", err)
	}

	return id, nil

}
