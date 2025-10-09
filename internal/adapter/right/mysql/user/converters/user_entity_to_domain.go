package userconverters

import (
	"fmt"
	"time"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

func UserEntityToDomain(entity []any) (user usermodel.UserInterface, err error) {
	user = usermodel.NewUser()

	id, ok := entity[0].(int64)
	if !ok {
		return nil, fmt.Errorf("user converter: invalid id type %T", entity[0])
	}
	user.SetID(id)

	full_name, ok := entity[1].([]byte)
	if !ok {
		return nil, fmt.Errorf("user converter: invalid full_name type %T", entity[1])
	}
	user.SetFullName(string(full_name))

	if entity[2] != nil {
		nick_name, ok := entity[2].([]byte)
		if !ok {
			return nil, fmt.Errorf("user converter: invalid nick_name type %T", entity[2])
		}
		user.SetNickName(string(nick_name))
	}

	national_id, ok := entity[3].([]byte)
	if !ok {
		return nil, fmt.Errorf("user converter: invalid national_id type %T", entity[3])
	}
	user.SetNationalID(string(national_id))

	if entity[4] != nil {
		creci_number, ok := entity[4].([]byte)
		if !ok {
			return nil, fmt.Errorf("user converter: invalid creci_number type %T", entity[4])
		}
		user.SetCreciNumber(string(creci_number))
	}

	if entity[5] != nil {
		creci_state, ok := entity[5].([]byte)
		if !ok {
			return nil, fmt.Errorf("user converter: invalid creci_state type %T", entity[5])
		}
		user.SetCreciState(string(creci_state))
	}

	if entity[6] != nil {
		creci_validity, ok := entity[6].(time.Time)
		if !ok {
			return nil, fmt.Errorf("user converter: invalid creci_validity type %T", entity[6])
		}
		user.SetCreciValidity(creci_validity)
	}

	born_at, ok := entity[7].(time.Time)
	if !ok {
		return nil, fmt.Errorf("user converter: invalid born_at type %T", entity[7])
	}
	user.SetBornAt(born_at)

	phone_number, ok := entity[8].([]byte)
	if !ok {
		return nil, fmt.Errorf("user converter: invalid phone_number type %T", entity[8])
	}
	user.SetPhoneNumber(string(phone_number))

	email, ok := entity[9].([]byte)
	if !ok {
		return nil, fmt.Errorf("user converter: invalid email type %T", entity[9])
	}
	user.SetEmail(string(email))

	zip_code, ok := entity[10].([]byte)
	if !ok {
		return nil, fmt.Errorf("user converter: invalid zip_code type %T", entity[10])
	}
	user.SetZipCode(string(zip_code))

	street, ok := entity[11].([]byte)
	if !ok {
		return nil, fmt.Errorf("user converter: invalid street type %T", entity[11])
	}
	user.SetStreet(string(street))

	number, ok := entity[12].([]byte)
	if !ok {
		return nil, fmt.Errorf("user converter: invalid number type %T", entity[12])
	}
	user.SetNumber(string(number))

	if entity[13] != nil {
		complement, ok := entity[13].([]byte)
		if !ok {
			return nil, fmt.Errorf("user converter: invalid complement type %T", entity[13])
		}
		user.SetComplement(string(complement))
	}

	neighborhood, ok := entity[14].([]byte)
	if !ok {
		return nil, fmt.Errorf("user converter: invalid neighborhood type %T", entity[14])
	}
	user.SetNeighborhood(string(neighborhood))

	city, ok := entity[15].([]byte)
	if !ok {
		return nil, fmt.Errorf("user converter: invalid city type %T", entity[15])
	}
	user.SetCity(string(city))

	state, ok := entity[16].([]byte)
	if !ok {
		return nil, fmt.Errorf("user converter: invalid state type %T", entity[16])
	}
	user.SetState(string(state))

	password, ok := entity[17].([]byte)
	if !ok {
		return nil, fmt.Errorf("user converter: invalid password type %T", entity[17])
	}
	user.SetPassword(string(password))

	// opt_status
	opt_status, ok := entity[18].(int64)
	if !ok {
		return nil, fmt.Errorf("user converter: invalid opt_status type %T", entity[18])
	}
	user.SetOptStatus(opt_status == 1)

	last_activity_at, ok := entity[19].(time.Time)
	if !ok {
		return nil, fmt.Errorf("user converter: invalid last_activity_at type %T", entity[19])
	}
	user.SetLastActivityAt(last_activity_at)

	deleted, ok := entity[20].(int64)
	if !ok {
		return nil, fmt.Errorf("user converter: invalid deleted type %T", entity[20])
	}
	user.SetDeleted(deleted == 1)

	if entity[21] != nil {
		last_sigin_attempt_at, ok := entity[21].(time.Time)
		if !ok {
			return nil, fmt.Errorf("user converter: invalid last_signin_attempt type %T", entity[21])
		}
		user.SetLastSignInAttempt(last_sigin_attempt_at)
	}

	return
}
