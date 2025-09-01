package userconverters

import (
	"fmt"
	"log/slog"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

func AgencyInviteEntityToDomain(entity []any) (domain usermodel.InviteInterface, err error) {
	domain = usermodel.NewInvite()

	id, ok := entity[0].(int64)
	if !ok {
		slog.Error("Error converting ID to int64", "value", entity[0])
		return nil, fmt.Errorf("agency invite: invalid id type %T", entity[0])
	}
	domain.SetID(id)

	agency_id, ok := entity[1].(int64)
	if !ok {
		slog.Error("Error converting ID to int64", "value", entity[1])
		return nil, fmt.Errorf("agency invite: invalid agency_id type %T", entity[1])
	}

	domain.SetAgencyID(agency_id)

	phone_number, ok := entity[2].([]byte)
	if !ok {
		slog.Error("Error converting status_reason to []byte", "value", entity[2])
		return nil, fmt.Errorf("agency invite: invalid phone_number type %T", entity[2])
	}
	domain.SetPhoneNumber(string(phone_number))

	return
}
