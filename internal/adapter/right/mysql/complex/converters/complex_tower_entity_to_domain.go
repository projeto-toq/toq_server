package complexrepoconverters

import (
	"log/slog"

	complexmodel "github.com/giulio-alfieri/toq_server/internal/core/model/complex_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func ComplexTowerEntityToDomain(entity []any) (complexTower complexmodel.ComplexTowerInterface, err error) {

	complexTower = complexmodel.NewComplexTower()

	id, ok := entity[0].(int64)
	if !ok {
		slog.Error("Error converting ID to int64", "ID", entity[0])
		return nil, utils.ErrInternalServer
	}
	complexTower.SetID(id)

	complex_id, ok := entity[1].(int64)
	if !ok {
		slog.Error("Error converting complex_id to int64", "complex_id", entity[1])
		return nil, utils.ErrInternalServer
	}
	complexTower.SetComplexID(complex_id)

	tower, ok := entity[2].([]byte)
	if !ok {
		slog.Error("Error converting tower to []byte", "tower", entity[2])
		return nil, utils.ErrInternalServer
	}
	complexTower.SetTower(string(tower))

	if entity[3] == nil {
		floors, ok := entity[3].(int32)
		if !ok {
			slog.Error("Error converting floors to int32", "floors", entity[3])
			return nil, utils.ErrInternalServer
		}
		complexTower.SetFloors(int(floors))
	}

	if entity[4] == nil {
		total_units, ok := entity[4].(int32)
		if !ok {
			slog.Error("Error converting total_units to int32", "total_units", entity[4])
			return nil, utils.ErrInternalServer
		}
		complexTower.SetTotalUnits(int(total_units))
	}

	if entity[5] == nil {
		units_per_floor, ok := entity[5].(int32)
		if !ok {
			slog.Error("Error converting units_per_floor to int32", "units_per_floor", entity[5])
			return nil, utils.ErrInternalServer
		}
		complexTower.SetUnitsPerFloor(int(units_per_floor))
	}

	return
}
