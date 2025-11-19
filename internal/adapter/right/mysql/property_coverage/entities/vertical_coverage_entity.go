package propertycoverageentities

// VerticalCoverageEntity maps the vertical_complexes table joined with its lookup fields.
type VerticalCoverageEntity struct {
	ID                   int64
	Name                 string
	MainRegistration     string
	PropertyTypesBitmask int64
}
