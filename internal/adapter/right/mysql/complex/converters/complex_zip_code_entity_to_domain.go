package complexrepoconverters

import (
	"log/slog"

	complexmodel "github.com/giulio-alfieri/toq_server/internal/core/model/complex_model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ComplexZipCodeEntityToDomain(entity []any) (complexZipCode complexmodel.ComplexZipCodeInterface, err error) {

	complexZipCode = complexmodel.NewComplexZipCode()

	id, ok := entity[0].(int64)
	if !ok {
		slog.Error("Error converting ID to int64", "ID", entity[0])
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	complexZipCode.SetID(id)

	complex_id, ok := entity[1].(int64)
	if !ok {
		slog.Error("Error converting complex_id to int64", "complex_id", entity[1])
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	complexZipCode.SetComplexID(complex_id)

	zip_code, ok := entity[2].([]byte)
	if !ok {
		slog.Error("Error converting zip_code to []byte", "zip_code", entity[2])
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	complexZipCode.SetZipCode(string(zip_code))

	return
}
