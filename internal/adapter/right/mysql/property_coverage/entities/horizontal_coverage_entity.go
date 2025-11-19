package propertycoverageentities

// HorizontalCoverageEntity maps horizontal complexes with their allowed property types.
type HorizontalCoverageEntity struct {
	ID                   int64
	Name                 string
	MainRegistration     string
	PropertyTypesBitmask int64
}
