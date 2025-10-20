package complexrepoconverters

import (
	"fmt"
	"strconv"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
)

func ComplexSizeEntityToDomain(entity []any) (complexSize complexmodel.ComplexSizeInterface, err error) {

	complexSize = complexmodel.NewComplexSize()

	id, ok := entity[0].(int64)
	if !ok {
		return nil, fmt.Errorf("convert id to int64: %T", entity[0])
	}
	complexSize.SetID(id)

	complex_id, ok := entity[1].(int64)
	if !ok {
		return nil, fmt.Errorf("convert complex_id to int64: %T", entity[1])
	}
	complexSize.SetComplexID(complex_id)

	switch v := entity[2].(type) {
	case float64:
		complexSize.SetSize(v)
	case float32:
		complexSize.SetSize(float64(v))
	case []byte:
		parsed, errParse := strconv.ParseFloat(string(v), 64)
		if errParse != nil {
			return nil, fmt.Errorf("convert size to float64: %w", errParse)
		}
		complexSize.SetSize(parsed)
	default:
		return nil, fmt.Errorf("convert size to float64: %T", entity[2])
	}

	if len(entity) > 3 && entity[3] != nil {
		descriptionBytes, ok := entity[3].([]byte)
		if !ok {
			return nil, fmt.Errorf("convert description to []byte: %T", entity[3])
		}
		complexSize.SetDescription(string(descriptionBytes))
	}

	return
}
