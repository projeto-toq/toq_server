package propertycoverageentities

import (
	"database/sql"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
)

// ManagedComplexEntity agrega colunas comuns retornadas nas consultas admin de complexos.
// Campos nullable usam sql.Null* para refletir o schema (street/number/neighborhood/reception_phone/main_registration).
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

// VerticalComplexTowerEntity espelha vertical_complex_towers (floors/total_units/units_per_floor são NOT NULL com default 0).
type VerticalComplexTowerEntity struct {
	ID                int64
	VerticalComplexID int64
	Tower             string
	Floors            int
	TotalUnits        int
	UnitsPerFloor     int
}

// VerticalComplexSizeEntity espelha vertical_complex_sizes (description é nullable).
type VerticalComplexSizeEntity struct {
	ID                int64
	VerticalComplexID int64
	Size              float64
	Description       sql.NullString
}

// HorizontalComplexZipCodeEntity espelha horizontal_complex_zip_codes.
type HorizontalComplexZipCodeEntity struct {
	ID                  int64
	HorizontalComplexID int64
	ZipCode             string
}
