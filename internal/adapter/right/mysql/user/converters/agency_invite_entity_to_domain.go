package userconverters

import (
	"log/slog"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func AgencyInviteEntityToDomain(entity []any) (domain usermodel.InviteInterface, err error) {
	domain = usermodel.NewInvite()

	id, ok := entity[0].(int64)
	if !ok {
		slog.Error("Error converting ID to int64", "value", entity[0])
		return nil, utils.ErrInternalServer
	}
	domain.SetID(id)

	agency_id, ok := entity[1].(int64)
	if !ok {
		slog.Error("Error converting ID to int64", "value", entity[1])
		return nil, utils.ErrInternalServer
	}

	domain.SetAgencyID(agency_id)

	phone_number, ok := entity[2].([]byte)
	if !ok {
		slog.Error("Error converting status_reason to []byte", "value", entity[2])
		return nil, utils.ErrInternalServer
	}
	domain.SetPhoneNumber(string(phone_number))

	return
}
