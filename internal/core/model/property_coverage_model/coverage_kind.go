package propertycoveragemodel

// CoverageKind identifies the managed coverage entity type handled by admin endpoints.
type CoverageKind string

const (
	CoverageKindVertical   CoverageKind = "VERTICAL"
	CoverageKindHorizontal CoverageKind = "HORIZONTAL"
	CoverageKindStandalone CoverageKind = "STANDALONE"
)

// Valid reports whether the provided kind matches a supported constant.
func (k CoverageKind) Valid() bool {
	switch k {
	case CoverageKindVertical, CoverageKindHorizontal, CoverageKindStandalone:
		return true
	default:
		return false
	}
}
