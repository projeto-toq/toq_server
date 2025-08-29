package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"

	userconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) GetWrongSigninByUserID(ctx context.Context, tx *sql.Tx, id int64) (wrongSignin usermodel.WrongSigninInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT * FROM temp_wrong_signin WHERE user_id = ?;`

	entities, err := ua.Read(ctx, tx, query, id)
	if err != nil {
		slog.Error("mysqluseradapter/GetWrongSigninByUserID: error executing Read", "error", err)
		return nil, utils.ErrInternalServer
	}

	if len(entities) == 0 {
		return nil, utils.ErrInternalServer
	}

	if len(entities) > 1 {
		slog.Error("mysqluseradapter/GetWrongSigninByUserID: multiple roles found with the same role", "role", id)
		return nil, utils.ErrInternalServer
	}

	wrongSignin, err = userconverters.WrongSignInEntityToDomain(entities[0])
	if err != nil {
		return
	}

	return

}
