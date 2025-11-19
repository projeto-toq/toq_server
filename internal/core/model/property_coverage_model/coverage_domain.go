package propertycoveragemodel

import (
	"strings"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
)

// coverage implements CoverageInterface.
type coverage struct {
	complexName      string
	mainRegistration string
	propertyTypes    globalmodel.PropertyType
	source           CoverageSource
}

// NewCoverage creates an empty domain object ready to be filled by repositories.
func NewCoverage() CoverageInterface {
	return &coverage{}
}

func (c *coverage) ComplexName() string {
	return c.complexName
}

func (c *coverage) SetComplexName(name string) {
	c.complexName = strings.TrimSpace(name)
}

func (c *coverage) PropertyTypes() globalmodel.PropertyType {
	return c.propertyTypes
}

func (c *coverage) SetPropertyTypes(value globalmodel.PropertyType) {
	c.propertyTypes = value
}

func (c *coverage) Source() CoverageSource {
	return c.source
}

func (c *coverage) SetSource(source CoverageSource) {
	c.source = source
}

func (c *coverage) MainRegistration() string {
	return c.mainRegistration
}

func (c *coverage) SetMainRegistration(value string) {
	c.mainRegistration = strings.TrimSpace(value)
}

func (c *coverage) HasComplex() bool {
	return c.complexName != ""
}

func (c *coverage) ToOutput() ResolvePropertyTypesOutput {
	return ResolvePropertyTypesOutput{
		PropertyTypes:    c.propertyTypes,
		ComplexName:      c.complexName,
		MainRegistration: c.mainRegistration,
		Source:           c.source,
	}
}
