package propertycoveragerepository

import (
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
)

// ListManagedComplexesParams configures paging and filtering for admin listings.
type ListManagedComplexesParams struct {
	Name         string
	ZipCode      string
	Number       string
	City         string
	State        string
	Sector       *propertycoveragemodel.Sector
	PropertyType *globalmodel.PropertyType
	Kind         *propertycoveragemodel.CoverageKind
	Limit        int
	Offset       int
}

// ListVerticalComplexTowersParams filters tower queries.
type ListVerticalComplexTowersParams struct {
	VerticalComplexID int64
	Tower             string
	Limit             int
	Offset            int
}

// ListVerticalComplexSizesParams filters size queries.
type ListVerticalComplexSizesParams struct {
	VerticalComplexID int64
	Limit             int
	Offset            int
}

// ListHorizontalComplexZipCodesParams filters horizontal CEP queries.
type ListHorizontalComplexZipCodesParams struct {
	HorizontalComplexID int64
	ZipCode             string
	Limit               int
	Offset              int
}
