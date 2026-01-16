package dto

// ComplexResponse representa um empreendimento e seus dados relacionados.
type ComplexResponse struct {
	CoverageType     string                   `json:"coverageType"`
	ID               int64                    `json:"id"`
	Name             string                   `json:"name"`
	ZipCode          string                   `json:"zipCode"`
	Street           string                   `json:"street"`
	Number           string                   `json:"number"`
	Neighborhood     string                   `json:"neighborhood"`
	City             string                   `json:"city"`
	State            string                   `json:"state"`
	PhoneNumber      string                   `json:"phoneNumber"`
	Sector           uint8                    `json:"sector"`
	MainRegistration string                   `json:"mainRegistration"`
	PropertyType     uint16                   `json:"propertyType"`
	Sizes            []ComplexSizeResponse    `json:"sizes,omitempty"`
	Towers           []ComplexTowerResponse   `json:"towers,omitempty"`
	ZipCodes         []ComplexZipCodeResponse `json:"zipCodes,omitempty"`
}

// ComplexTowerResponse descreve uma torre vinculada a um empreendimento.
type ComplexTowerResponse struct {
	ID            int64  `json:"id"`
	ComplexID     int64  `json:"complexId"`
	Tower         string `json:"tower"`
	Floors        *int   `json:"floors,omitempty"`
	TotalUnits    *int   `json:"totalUnits,omitempty"`
	UnitsPerFloor *int   `json:"unitsPerFloor,omitempty"`
}

// ComplexSizeResponse descreve um tamanho disponível de empreendimento.
type ComplexSizeResponse struct {
	ID          int64   `json:"id"`
	ComplexID   int64   `json:"complexId"`
	Size        float64 `json:"size"`
	Description string  `json:"description"`
}

// ComplexZipCodeResponse descreve um CEP vinculado a um empreendimento.
type ComplexZipCodeResponse struct {
	ID        int64  `json:"id"`
	ComplexID int64  `json:"complexId"`
	ZipCode   string `json:"zipCode"`
}

// GetComplexByAddressQuery representa os parâmetros de consulta para o endpoint público de complexo.
type GetComplexByAddressQuery struct {
	ZipCode string `form:"zipCode" binding:"required"`
	Number  string `form:"number"`
}

// AdminCreateComplexRequest representa o payload para criação de empreendimento.
type AdminCreateComplexRequest struct {
	CoverageType     string  `json:"coverageType" binding:"required,oneof=VERTICAL HORIZONTAL STANDALONE"`
	Name             string  `json:"name" binding:"required"`
	ZipCode          string  `json:"zipCode" binding:"required"`
	Street           string  `json:"street"`
	Number           string  `json:"number" binding:"required"`
	Neighborhood     string  `json:"neighborhood"`
	City             string  `json:"city" binding:"required"`
	State            string  `json:"state" binding:"required"`
	PhoneNumber      string  `json:"phoneNumber"`
	Sector           *uint8  `json:"sector" binding:"required"`
	MainRegistration string  `json:"mainRegistration"`
	PropertyType     *uint16 `json:"propertyType" binding:"required"`
}

// AdminUpdateComplexRequest representa o payload para atualização de empreendimento.
type AdminUpdateComplexRequest struct {
	ID               int64   `json:"id" binding:"required,min=1"`
	CoverageType     string  `json:"coverageType" binding:"required,oneof=VERTICAL HORIZONTAL STANDALONE"`
	Name             string  `json:"name" binding:"required"`
	ZipCode          string  `json:"zipCode" binding:"required"`
	Street           string  `json:"street"`
	Number           string  `json:"number" binding:"required"`
	Neighborhood     string  `json:"neighborhood"`
	City             string  `json:"city" binding:"required"`
	State            string  `json:"state" binding:"required"`
	PhoneNumber      string  `json:"phoneNumber"`
	Sector           *uint8  `json:"sector" binding:"required"`
	MainRegistration string  `json:"mainRegistration"`
	PropertyType     *uint16 `json:"propertyType" binding:"required"`
}

// AdminDeleteComplexRequest representa o payload para exclusão de empreendimento.
type AdminDeleteComplexRequest struct {
	ID           int64  `json:"id" binding:"required,min=1"`
	CoverageType string `json:"coverageType" binding:"required,oneof=VERTICAL HORIZONTAL STANDALONE"`
}

// AdminListComplexesRequest representa filtros para listagem de empreendimentos.
type AdminListComplexesRequest struct {
	Name         string  `json:"name" form:"name"`
	ZipCode      string  `json:"zipCode" form:"zipCode"`
	Number       string  `json:"number" form:"number"`
	City         string  `json:"city" form:"city"`
	State        string  `json:"state" form:"state"`
	Sector       *uint8  `json:"sector" form:"sector"`
	PropertyType *uint16 `json:"propertyType" form:"propertyType"`
	CoverageType string  `json:"coverageType" form:"coverageType"`
	Page         int     `json:"page" form:"page"`
	Limit        int     `json:"limit" form:"limit"`
}

// AdminGetComplexDetailRequest representa o payload para obter detalhes de um empreendimento.
type AdminGetComplexDetailRequest struct {
	ID           int64  `json:"id" binding:"required,min=1"`
	CoverageType string `json:"coverageType" binding:"required,oneof=VERTICAL HORIZONTAL STANDALONE"`
}

// AdminListComplexesResponse encapsula a lista de empreendimentos retornada.
type AdminListComplexesResponse struct {
	Complexes []ComplexResponse `json:"complexes"`
	Page      int               `json:"page"`
	Limit     int               `json:"limit"`
}

// AdminCreateComplexTowerRequest representa o payload para criação de torre.
type AdminCreateComplexTowerRequest struct {
	ComplexID     int64  `json:"complexId" binding:"required,min=1"`
	Tower         string `json:"tower" binding:"required"`
	Floors        *int   `json:"floors"`
	TotalUnits    *int   `json:"totalUnits"`
	UnitsPerFloor *int   `json:"unitsPerFloor"`
}

// AdminUpdateComplexTowerRequest representa o payload para atualização de torre.
type AdminUpdateComplexTowerRequest struct {
	ID            int64  `json:"id" binding:"required,min=1"`
	ComplexID     int64  `json:"complexId" binding:"required,min=1"`
	Tower         string `json:"tower" binding:"required"`
	Floors        *int   `json:"floors"`
	TotalUnits    *int   `json:"totalUnits"`
	UnitsPerFloor *int   `json:"unitsPerFloor"`
}

// AdminDeleteComplexTowerRequest representa o payload para exclusão de torre.
type AdminDeleteComplexTowerRequest struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

// AdminListComplexTowersRequest representa filtros para listagem de torres.
type AdminListComplexTowersRequest struct {
	ComplexID int64  `json:"complexId" form:"complexId"`
	Tower     string `json:"tower" form:"tower"`
	Page      int    `json:"page" form:"page"`
	Limit     int    `json:"limit" form:"limit"`
}

// AdminListComplexTowersResponse encapsula a lista de torres retornada.
type AdminListComplexTowersResponse struct {
	Towers []ComplexTowerResponse `json:"towers"`
	Page   int                    `json:"page"`
	Limit  int                    `json:"limit"`
}

// AdminCreateComplexSizeRequest representa o payload para criação de tamanho.
type AdminCreateComplexSizeRequest struct {
	ComplexID   int64   `json:"complexId" binding:"required,min=1"`
	Size        float64 `json:"size" binding:"required"`
	Description string  `json:"description"`
}

// AdminUpdateComplexSizeRequest representa o payload para atualização de tamanho.
type AdminUpdateComplexSizeRequest struct {
	ID          int64   `json:"id" binding:"required,min=1"`
	ComplexID   int64   `json:"complexId" binding:"required,min=1"`
	Size        float64 `json:"size" binding:"required"`
	Description string  `json:"description"`
}

// AdminDeleteComplexSizeRequest representa o payload para exclusão de tamanho.
type AdminDeleteComplexSizeRequest struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

// AdminListComplexSizesRequest representa filtros para tamanhos.
type AdminListComplexSizesRequest struct {
	ComplexID int64 `json:"complexId" form:"complexId"`
	Page      int   `json:"page" form:"page"`
	Limit     int   `json:"limit" form:"limit"`
}

// AdminGetComplexSizeDetailRequest represents the payload to fetch a complex size detail.
type AdminGetComplexSizeDetailRequest struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

// AdminGetComplexTowerDetailRequest represents the payload to fetch a complex tower detail.
type AdminGetComplexTowerDetailRequest struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

// AdminListComplexSizesResponse encapsula tamanhos retornados.
type AdminListComplexSizesResponse struct {
	Sizes []ComplexSizeResponse `json:"sizes"`
	Page  int                   `json:"page"`
	Limit int                   `json:"limit"`
}

// AdminCreateComplexZipCodeRequest representa o payload para criação de CEP.
type AdminCreateComplexZipCodeRequest struct {
	ComplexID int64  `json:"complexId" binding:"required,min=1"`
	ZipCode   string `json:"zipCode" binding:"required"`
}

// AdminUpdateComplexZipCodeRequest representa o payload para atualização de CEP.
type AdminUpdateComplexZipCodeRequest struct {
	ID        int64  `json:"id" binding:"required,min=1"`
	ComplexID int64  `json:"complexId" binding:"required,min=1"`
	ZipCode   string `json:"zipCode" binding:"required"`
}

// AdminDeleteComplexZipCodeRequest representa o payload para exclusão de CEP.
type AdminDeleteComplexZipCodeRequest struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

// AdminListComplexZipCodesRequest representa filtros para CEPs vinculados a empreendimentos.
type AdminListComplexZipCodesRequest struct {
	ComplexID int64  `json:"complexId" form:"complexId"`
	ZipCode   string `json:"zipCode" form:"zipCode"`
	Page      int    `json:"page" form:"page"`
	Limit     int    `json:"limit" form:"limit"`
}

// AdminListComplexZipCodesResponse encapsula os CEPs retornados.
type AdminListComplexZipCodesResponse struct {
	ZipCodes []ComplexZipCodeResponse `json:"zipCodes"`
	Page     int                      `json:"page"`
	Limit    int                      `json:"limit"`
}

// AdminGetComplexZipCodeDetailRequest represents the payload to fetch a complex zip code detail.
type AdminGetComplexZipCodeDetailRequest struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

// ListingComplexItemResponse is a minimal payload for public listing flows.
type ListingComplexItemResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
