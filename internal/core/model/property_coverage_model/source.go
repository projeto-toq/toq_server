package propertycoveragemodel

// CoverageSource identifies which dataset matched the request.
type CoverageSource string

const (
	CoverageSourceVertical   CoverageSource = "vertical_complex"
	CoverageSourceHorizontal CoverageSource = "horizontal_complex"
	CoverageSourceStandalone CoverageSource = "no_complex"
)
