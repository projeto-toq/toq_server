package propertycoveragemodel

// Sector mirrors the segment classification previously available in the legacy complex model.
type Sector uint8

const (
	SectorResidential Sector = iota
	SectorCommercial
	SectorMixed
)
