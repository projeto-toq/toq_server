package propertycoverageentities

// NoComplexCoverageEntity armazena tipos permitidos para áreas standalone (no_complex_zip_codes).
// zip_code é a chave de lookup; PropertyTypesBitmask segue schema INT UNSIGNED.
type NoComplexCoverageEntity struct {
	ZipCode              string
	PropertyTypesBitmask int64
}
