package complexrepoconverters

import (
	"log/slog"

	complexmodel "github.com/giulio-alfieri/toq_server/internal/core/model/complex_model"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func ComplexEntityToDomain(entity []any) (complex complexmodel.ComplexInterface, err error) {

	complex = complexmodel.NewComplex()

	id, ok := entity[0].(int64)
	if !ok {
		slog.Error("Error converting ID to int64", "ID", entity[0])
		return nil, utils.ErrInternalServer
	}
	complex.SetID(id)

	name, ok := entity[1].([]byte)
	if !ok {
		slog.Error("Error converting name to []byte", "name", entity[1])
		return nil, utils.ErrInternalServer
	}
	complex.SetName(string(name))

	zip_code, ok := entity[2].([]byte)
	if !ok {
		slog.Error("Error converting zip_code to []byte", "zip_code", entity[2])
		return nil, utils.ErrInternalServer
	}
	complex.SetZipCode(string(zip_code))

	if entity[3] == nil {
		street, ok := entity[3].([]byte)
		if !ok {
			slog.Error("Error converting street to []byte", "street", entity[3])
			return nil, utils.ErrInternalServer
		}
		complex.SetStreet(string(street))
	}

	number, ok := entity[4].([]byte)
	if !ok {
		slog.Error("Error converting number to []byte", "number", entity[4])
		return nil, utils.ErrInternalServer
	}
	complex.SetNumber(string(number))

	if entity[5] == nil {
		neighborhood, ok := entity[5].([]byte)
		if !ok {
			slog.Error("Error converting neighborhood to []byte", "neighborhood", entity[5])
			return nil, utils.ErrInternalServer
		}
		complex.SetNeighborhood(string(neighborhood))
	}

	city, ok := entity[6].([]byte)
	if !ok {
		slog.Error("Error converting city to []byte", "city", entity[6])
		return nil, utils.ErrInternalServer
	}
	complex.SetCity(string(city))

	state, ok := entity[7].([]byte)
	if !ok {
		slog.Error("Error converting state to []byte", "state", entity[7])
		return nil, utils.ErrInternalServer
	}
	complex.SetState(string(state))

	if entity[8] == nil {
		reception_phone, ok := entity[8].([]byte)
		if !ok {
			slog.Error("Error converting reception_phone to []byte", "reception_phone", entity[8])
			return nil, utils.ErrInternalServer
		}
		complex.SetPhoneNumber(string(reception_phone))
	}

	sector, ok := entity[9].(int64)
	if !ok {
		slog.Error("Error converting sector to int64", "sector", entity[9])
		return nil, utils.ErrInternalServer
	}
	complex.SetSector(complexmodel.Sector(sector))

	if entity[10] == nil {
		main_registration, ok := entity[10].([]byte)
		if !ok {
			slog.Error("Error converting main_registration to []byte", "main_registration", entity[10])
			return nil, utils.ErrInternalServer
		}
		complex.SetMainRegistration(string(main_registration))
	}

	property_type, ok := entity[11].(int64)
	if !ok {
		slog.Error("Error converting property_type to int64", "property_type", entity[11])
		return nil, utils.ErrInternalServer
	}
	complex.SetPropertyType(globalmodel.PropertyType(property_type))

	return
}
