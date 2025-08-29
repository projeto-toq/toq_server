package complexrepoconverters

import (
	"log/slog"

	complexmodel "github.com/giulio-alfieri/toq_server/internal/core/model/complex_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func ComplexZipCodeEntityToDomain(entity []any) (complexZipCode complexmodel.ComplexZipCodeInterface, err error) {

	complexZipCode = complexmodel.NewComplexZipCode()

	id, ok := entity[0].(int64)
	if !ok {
		slog.Error("Error converting ID to int64", "ID", entity[0])
		return nil, utils.ErrInternalServer
	}
	complexZipCode.SetID(id)

	complex_id, ok := entity[1].(int64)
	if !ok {
		slog.Error("Error converting complex_id to int64", "complex_id", entity[1])
		return nil, utils.ErrInternalServer
	}
	complexZipCode.SetComplexID(complex_id)

	zip_code, ok := entity[2].([]byte)
	if !ok {
		slog.Error("Error converting zip_code to []byte", "zip_code", entity[2])
		return nil, utils.ErrInternalServer
	}
	complexZipCode.SetZipCode(string(zip_code))

	return
}
