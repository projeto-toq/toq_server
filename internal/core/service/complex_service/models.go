package complexservices

import (
	"strings"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
)

// CreateComplexInput define os campos necessários para criar um empreendimento.
type CreateComplexInput struct {
	Name             string
	ZipCode          string
	Street           string
	Number           string
	Neighborhood     string
	City             string
	State            string
	PhoneNumber      string
	Sector           complexmodel.Sector
	MainRegistration string
	PropertyType     globalmodel.PropertyType
}

// UpdateComplexInput define os campos para atualizar um empreendimento existente.
type UpdateComplexInput struct {
	ID               int64
	Name             string
	ZipCode          string
	Street           string
	Number           string
	Neighborhood     string
	City             string
	State            string
	PhoneNumber      string
	Sector           complexmodel.Sector
	MainRegistration string
	PropertyType     globalmodel.PropertyType
}

// ListComplexesInput contém filtros opcionais para listagem de empreendimentos.
type ListComplexesInput struct {
	Name         string
	ZipCode      string
	City         string
	State        string
	Sector       *complexmodel.Sector
	PropertyType *globalmodel.PropertyType
	Page         int
	Limit        int
}

// CreateComplexTowerInput representa os dados para criar uma torre.
type CreateComplexTowerInput struct {
	ComplexID     int64
	Tower         string
	Floors        *int
	TotalUnits    *int
	UnitsPerFloor *int
}

// UpdateComplexTowerInput representa os dados para atualizar uma torre existente.
type UpdateComplexTowerInput struct {
	ID            int64
	ComplexID     int64
	Tower         string
	Floors        *int
	TotalUnits    *int
	UnitsPerFloor *int
}

// ListComplexTowersInput define filtros básicos para listar torres.
type ListComplexTowersInput struct {
	ComplexID int64
	Tower     string
	Page      int
	Limit     int
}

// CreateComplexSizeInput representa os dados para criar um tamanho disponível.
type CreateComplexSizeInput struct {
	ComplexID   int64
	Size        float64
	Description string
}

// UpdateComplexSizeInput representa os dados para atualizar um tamanho existente.
type UpdateComplexSizeInput struct {
	ID          int64
	ComplexID   int64
	Size        float64
	Description string
}

// ListComplexSizesInput define filtros para listar tamanhos de um empreendimento.
type ListComplexSizesInput struct {
	ComplexID int64
	Page      int
	Limit     int
}

// CreateComplexZipCodeInput representa os dados para criar um CEP associado ao empreendimento.
type CreateComplexZipCodeInput struct {
	ComplexID int64
	ZipCode   string
}

// UpdateComplexZipCodeInput representa os dados para atualizar um CEP associado.
type UpdateComplexZipCodeInput struct {
	ID        int64
	ComplexID int64
	ZipCode   string
}

// ListComplexZipCodesInput define filtros para listar CEPs associados a um empreendimento.
type ListComplexZipCodesInput struct {
	ComplexID int64
	ZipCode   string
	Page      int
	Limit     int
}

// ListSizesByAddressInput contém os parâmetros para consultar tamanhos a partir de CEP e número.
type ListSizesByAddressInput struct {
	ZipCode string
	Number  string
}

// sanitizeString remove espaços extras, retornando string vazia quando apropriado.
func sanitizeString(value string) string {
	return strings.TrimSpace(value)
}
