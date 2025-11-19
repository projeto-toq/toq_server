package propertycoverageservice

import (
	"strings"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
)

// CreateComplexInput contains normalized data to create a managed coverage entry.
type CreateComplexInput struct {
	Kind             propertycoveragemodel.CoverageKind
	Name             string
	ZipCode          string
	Street           string
	Number           string
	Neighborhood     string
	City             string
	State            string
	ReceptionPhone   string
	Sector           propertycoveragemodel.Sector
	MainRegistration string
	PropertyType     globalmodel.PropertyType
}

// UpdateComplexInput extends CreateComplexInput with the entity identifier.
type UpdateComplexInput struct {
	CreateComplexInput
	ID int64
}

// DeleteComplexInput identifies the coverage entry that must be removed.
type DeleteComplexInput struct {
	ID   int64
	Kind propertycoveragemodel.CoverageKind
}

// GetComplexDetailInput identifies the coverage entry to fetch.
type GetComplexDetailInput struct {
	ID   int64
	Kind propertycoveragemodel.CoverageKind
}

// ListComplexesInput configures filters for admin listings.
type ListComplexesInput struct {
	Name         string
	ZipCode      string
	Number       string
	City         string
	State        string
	Sector       *propertycoveragemodel.Sector
	PropertyType *globalmodel.PropertyType
	Kind         *propertycoveragemodel.CoverageKind
	Page         int
	Limit        int
}

// CreateComplexTowerInput holds tower payload data.
type CreateComplexTowerInput struct {
	VerticalComplexID int64
	Tower             string
	Floors            *int
	TotalUnits        *int
	UnitsPerFloor     *int
}

// UpdateComplexTowerInput extends the tower payload with the identifier.
type UpdateComplexTowerInput struct {
	CreateComplexTowerInput
	ID int64
}

// ListComplexTowersInput filters tower queries.
type ListComplexTowersInput struct {
	VerticalComplexID int64
	Tower             string
	Page              int
	Limit             int
}

// CreateComplexSizeInput holds size payload data.
type CreateComplexSizeInput struct {
	VerticalComplexID int64
	Size              float64
	Description       string
}

// UpdateComplexSizeInput extends the size payload with the identifier.
type UpdateComplexSizeInput struct {
	CreateComplexSizeInput
	ID int64
}

// ListComplexSizesInput filters size queries.
type ListComplexSizesInput struct {
	VerticalComplexID int64
	Page              int
	Limit             int
}

// CreateComplexZipCodeInput holds horizontal CEP payload data.
type CreateComplexZipCodeInput struct {
	HorizontalComplexID int64
	ZipCode             string
}

// UpdateComplexZipCodeInput extends the CEP payload with the identifier.
type UpdateComplexZipCodeInput struct {
	CreateComplexZipCodeInput
	ID int64
}

// ListComplexZipCodesInput filters horizontal CEP queries.
type ListComplexZipCodesInput struct {
	HorizontalComplexID int64
	ZipCode             string
	Page                int
	Limit               int
}

// GetComplexByAddressInput is reused by public endpoints to fetch complex details.
type GetComplexByAddressInput struct {
	ZipCode string
	Number  string
}

func sanitizeString(value string) string {
	return strings.TrimSpace(value)
}

func sanitizePagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return page, limit
}
