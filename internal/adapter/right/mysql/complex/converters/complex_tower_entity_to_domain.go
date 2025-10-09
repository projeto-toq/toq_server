package complexrepoconverters

import (
	"fmt"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
)

func ComplexTowerEntityToDomain(entity []any) (complexTower complexmodel.ComplexTowerInterface, err error) {

	complexTower = complexmodel.NewComplexTower()

	id, ok := entity[0].(int64)
	if !ok {
		return nil, fmt.Errorf("convert id to int64: %T", entity[0])
	}
	complexTower.SetID(id)

	complex_id, ok := entity[1].(int64)
	if !ok {
		return nil, fmt.Errorf("convert complex_id to int64: %T", entity[1])
	}
	complexTower.SetComplexID(complex_id)

	tower, ok := entity[2].([]byte)
	if !ok {
		return nil, fmt.Errorf("convert tower to []byte: %T", entity[2])
	}
	complexTower.SetTower(string(tower))

	if entity[3] != nil {
		floors, ok := entity[3].(int64)
		if !ok {
			return nil, fmt.Errorf("convert floors to int64: %T", entity[3])
		}
		complexTower.SetFloors(int(floors))
	}

	if entity[4] != nil {
		total_units, ok := entity[4].(int64)
		if !ok {
			return nil, fmt.Errorf("convert total_units to int64: %T", entity[4])
		}
		complexTower.SetTotalUnits(int(total_units))
	}

	if entity[5] != nil {
		units_per_floor, ok := entity[5].(int64)
		if !ok {
			return nil, fmt.Errorf("convert units_per_floor to int64: %T", entity[5])
		}
		complexTower.SetUnitsPerFloor(int(units_per_floor))
	}

	return
}
