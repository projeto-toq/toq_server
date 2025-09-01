package complexrepoconverters

import (
	"fmt"
	"log/slog"

	complexmodel "github.com/giulio-alfieri/toq_server/internal/core/model/complex_model"
)

func ComplexSizeEntityToDomain(entity []any) (complexSize complexmodel.ComplexSizeInterface, err error) {

	complexSize = complexmodel.NewComplexSize()

	id, ok := entity[0].(int64)
	if !ok {
		slog.Error("Error converting ID to int64", "ID", entity[0])
		return nil, fmt.Errorf("convert id to int64: %T", entity[0])
	}
	complexSize.SetID(id)

	complex_id, ok := entity[1].(int64)
	if !ok {
		slog.Error("Error converting complex_id to int64", "complex_id", entity[1])
		return nil, fmt.Errorf("convert complex_id to int64: %T", entity[1])
	}
	complexSize.SetComplexID(complex_id)

	size, ok := entity[2].(float64)
	if !ok {
		slog.Error("Error converting size to float64", "size", entity[2])
		return nil, fmt.Errorf("convert size to float64: %T", entity[2])
	}
	complexSize.SetSize(size)

	return
}
