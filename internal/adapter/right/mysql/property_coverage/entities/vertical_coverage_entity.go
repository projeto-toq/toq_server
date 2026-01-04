package propertycoverageentities

import "database/sql"

// VerticalCoverageEntity representa linha de vertical_complexes usada em consultas de cobertura.
// main_registration é nullable no schema, por isso sql.NullString é usado para leitura segura.
type VerticalCoverageEntity struct {
	ID                   int64
	Name                 string
	MainRegistration     sql.NullString
	PropertyTypesBitmask int64
}
