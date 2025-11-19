package propertycoverageentities

import (
	"database/sql"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
)

// ManagedComplexEntity aggregates the common columns returned by admin queries.
type ManagedComplexEntity struct {
	ID               int64
	Kind             propertycoveragemodel.CoverageKind
	Name             sql.NullString
	ZipCode          string
	Street           sql.NullString
	Number           sql.NullString
	Neighborhood     sql.NullString
	City             string
	State            string
	ReceptionPhone   sql.NullString
	Sector           uint8
	MainRegistration sql.NullString
	PropertyTypes    uint16
}

// VerticalComplexTowerEntity mirrors the vertical_complex_towers table.
type VerticalComplexTowerEntity struct {
	ID                int64
	VerticalComplexID int64
	Tower             string
	Floors            int
	TotalUnits        int
	UnitsPerFloor     int
}

// VerticalComplexSizeEntity mirrors the vertical_complex_sizes table.
type VerticalComplexSizeEntity struct {
	ID                int64
	VerticalComplexID int64
	Size              float64
	Description       sql.NullString
}

// HorizontalComplexZipCodeEntity mirrors the horizontal_complex_zip_codes table.
type HorizontalComplexZipCodeEntity struct {
	ID                  int64
	HorizontalComplexID int64
	ZipCode             string
}
