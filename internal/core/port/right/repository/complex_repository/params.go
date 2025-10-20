package complexrepository

import (
	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
)

// ListComplexesParams define filtros opcionais utilizados na consulta de empreendimentos.
type ListComplexesParams struct {
	Name         string
	ZipCode      string
	City         string
	State        string
	Sector       *complexmodel.Sector
	PropertyType *globalmodel.PropertyType
	Limit        int
	Offset       int
}

// ListComplexTowersParams define filtros opcionais para consulta de torres vinculadas a um empreendimento.
type ListComplexTowersParams struct {
	ComplexID int64
	Tower     string
	Limit     int
	Offset    int
}

// ListComplexSizesParams define filtros opcionais para consulta de tamanhos cadastrados.
type ListComplexSizesParams struct {
	ComplexID int64
	Limit     int
	Offset    int
}

// ListComplexZipCodesParams define filtros opcionais para consulta de CEPs vinculados.
type ListComplexZipCodesParams struct {
	ComplexID int64
	ZipCode   string
	Limit     int
	Offset    int
}
