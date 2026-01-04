package propertycoverageentities

import "database/sql"

// HorizontalCoverageEntity representa linha de horizontal_complexes usada em consultas de cobertura.
// main_registration é nullable na tabela, portanto sql.NullString é utilizado para evitar panics em Scan.
type HorizontalCoverageEntity struct {
	ID                   int64
	Name                 string
	MainRegistration     sql.NullString
	PropertyTypesBitmask int64
}
