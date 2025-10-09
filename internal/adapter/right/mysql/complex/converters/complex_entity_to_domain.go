package complexrepoconverters

import (
	"fmt"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
)

func ComplexEntityToDomain(entity []any) (complex complexmodel.ComplexInterface, err error) {

	complex = complexmodel.NewComplex()

	id, ok := entity[0].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid complex ID type: %T", entity[0])
	}
	complex.SetID(id)

	name, ok := entity[1].([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid name type: %T", entity[1])
	}
	complex.SetName(string(name))

	zip_code, ok := entity[2].([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid zip_code type: %T", entity[2])
	}
	complex.SetZipCode(string(zip_code))

	if entity[3] != nil {
		street, ok := entity[3].([]byte)
		if !ok {
			return nil, fmt.Errorf("invalid street type: %T", entity[3])
		}
		complex.SetStreet(string(street))
	}

	number, ok := entity[4].([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid number type: %T", entity[4])
	}
	complex.SetNumber(string(number))

	if entity[5] != nil {
		neighborhood, ok := entity[5].([]byte)
		if !ok {
			return nil, fmt.Errorf("invalid neighborhood type: %T", entity[5])
		}
		complex.SetNeighborhood(string(neighborhood))
	}

	city, ok := entity[6].([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid city type: %T", entity[6])
	}
	complex.SetCity(string(city))

	state, ok := entity[7].([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid state type: %T", entity[7])
	}
	complex.SetState(string(state))

	if entity[8] != nil {
		reception_phone, ok := entity[8].([]byte)
		if !ok {
			return nil, fmt.Errorf("invalid reception_phone type: %T", entity[8])
		}
		complex.SetPhoneNumber(string(reception_phone))
	}

	sector, ok := entity[9].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid sector type: %T", entity[9])
	}
	complex.SetSector(complexmodel.Sector(sector))

	if entity[10] != nil {
		main_registration, ok := entity[10].([]byte)
		if !ok {
			return nil, fmt.Errorf("invalid main_registration type: %T", entity[10])
		}
		complex.SetMainRegistration(string(main_registration))
	}

	property_type, ok := entity[11].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid property_type type: %T", entity[11])
	}
	complex.SetPropertyType(globalmodel.PropertyType(property_type))

	return
}
