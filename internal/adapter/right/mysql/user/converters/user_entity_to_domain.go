package userconverters

import (
	"log/slog"
	"time"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func UserEntityToDomain(entity []any) (user usermodel.UserInterface, err error) {
	user = usermodel.NewUser()

	id, ok := entity[0].(int64)
	if !ok {
		slog.Error("Error converting ID to int64", "ID", entity[0])
		return nil, utils.ErrInternalServer
	}
	user.SetID(id)

	full_name, ok := entity[1].([]byte)
	if !ok {
		slog.Error("Error converting full_name to []byte", "full_name", entity[1])
		return nil, utils.ErrInternalServer
	}
	user.SetFullName(string(full_name))

	if entity[2] != nil {
		nick_name, ok := entity[2].([]byte)
		if !ok {
			slog.Error("Error converting nick_name to string", "nick_name", entity[2])
			return nil, utils.ErrInternalServer
		}
		user.SetNickName(string(nick_name))
	}

	national_id, ok := entity[3].([]byte)
	if !ok {
		slog.Error("Error converting national_id to string", "national_id", entity[3])
		return nil, utils.ErrInternalServer
	}
	user.SetNationalID(string(national_id))

	if entity[4] != nil {
		creci_number, ok := entity[4].([]byte)
		if !ok {
			slog.Error("Error converting creci_number to string", "creci_number", entity[4])
			return nil, utils.ErrInternalServer
		}
		user.SetCreciNumber(string(creci_number))
	}

	if entity[5] != nil {
		creci_state, ok := entity[5].([]byte)
		if !ok {
			slog.Error("Error converting creci_state to string", "creci_state", entity[5])
			return nil, utils.ErrInternalServer
		}
		user.SetCreciState(string(creci_state))
	}

	if entity[6] != nil {
		creci_validity, ok := entity[6].(time.Time)
		if !ok {
			slog.Error("Error converting creci_validity to time.Time", "creci_validity", entity[6])
			return nil, utils.ErrInternalServer
		}
		user.SetCreciValidity(creci_validity)
	}

	born_at, ok := entity[7].(time.Time)
	if !ok {
		slog.Error("Error converting born_at to time.Time", "born_at", entity[7])
		return nil, utils.ErrInternalServer
	}
	user.SetBornAt(born_at)

	phone_number, ok := entity[8].([]byte)
	if !ok {
		slog.Error("Error converting phone_number to string", "phone_number", entity[8])
		return nil, utils.ErrInternalServer
	}
	user.SetPhoneNumber(string(phone_number))

	email, ok := entity[9].([]byte)
	if !ok {
		slog.Error("Error converting email to string", "email", entity[9])
		return nil, utils.ErrInternalServer
	}
	user.SetEmail(string(email))

	zip_code, ok := entity[10].([]byte)
	if !ok {
		slog.Error("Error converting zip_code to string", "zip_code", entity[10])
		return nil, utils.ErrInternalServer
	}
	user.SetZipCode(string(zip_code))

	street, ok := entity[11].([]byte)
	if !ok {
		slog.Error("Error converting street to string", "street", entity[11])
		return nil, utils.ErrInternalServer
	}
	user.SetStreet(string(street))

	number, ok := entity[12].([]byte)
	if !ok {
		slog.Error("Error converting number to string", "number", entity[12])
		return nil, utils.ErrInternalServer
	}
	user.SetNumber(string(number))

	if entity[13] != nil {
		complement, ok := entity[13].([]byte)
		if !ok {
			slog.Error("Error converting complement to string", "complement", entity[13])
			return nil, utils.ErrInternalServer
		}
		user.SetComplement(string(complement))
	}

	neighborhood, ok := entity[14].([]byte)
	if !ok {
		slog.Error("Error converting neighborhood to string", "neighborhood", entity[14])
		return nil, utils.ErrInternalServer
	}
	user.SetNeighborhood(string(neighborhood))

	city, ok := entity[15].([]byte)
	if !ok {
		slog.Error("Error converting city to string", "city", entity[15])
		return nil, utils.ErrInternalServer
	}
	user.SetCity(string(city))

	state, ok := entity[16].([]byte)
	if !ok {
		slog.Error("Error converting state to string", "state", entity[16])
		return nil, utils.ErrInternalServer
	}
	user.SetState(string(state))

	password, ok := entity[17].([]byte)
	if !ok {
		slog.Error("Error converting password to string", "password", entity[17])
		return nil, utils.ErrInternalServer
	}
	user.SetPassword(string(password))

	// opt_status
	opt_status, ok := entity[18].(int64)
	if !ok {
		slog.Error("Error converting opt_status to bool", "opt_status", entity[18])
		return nil, utils.ErrInternalServer
	}
	user.SetOptStatus(opt_status == 1)

	last_activity_at, ok := entity[19].(time.Time)
	if !ok {
		slog.Error("Error converting last_activity_at to time.Time", "last_activity_at", entity[19])
		return nil, utils.ErrInternalServer
	}
	user.SetLastActivityAt(last_activity_at)

	deleted, ok := entity[20].(int64)
	if !ok {
		slog.Error("Error converting deleted to bool", "deleted", entity[20])
		return nil, utils.ErrInternalServer
	}
	user.SetDeleted(deleted == 1)

	if entity[21] != nil {
		last_sigin_attempt_at, ok := entity[21].(time.Time)
		if !ok {
			slog.Error("Error converting last_sigin_attempt_at to time.Time", "last_sigin_attempt_at", entity[21])
			return nil, utils.ErrInternalServer
		}
		user.SetLastSignInAttempt(last_sigin_attempt_at)
	}

	return
}
