package propertycoveragemodel

import globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"

// CoverageInterface describes the minimum contract required to transport
// property coverage metadata across layers.
type CoverageInterface interface {
	ComplexName() string
	SetComplexName(string)
	MainRegistration() string
	SetMainRegistration(string)
	PropertyTypes() globalmodel.PropertyType
	SetPropertyTypes(globalmodel.PropertyType)
	Source() CoverageSource
	SetSource(CoverageSource)
	HasComplex() bool
	ToOutput() ResolvePropertyTypesOutput
}

// ResolvePropertyTypesInput holds the normalized parameters used to find
// property coverage information.
type ResolvePropertyTypesInput struct {
	ZipCode string
	Number  string
}

// ResolvePropertyTypesOutput mirrors the repository result in a service-friendly
// format so other domains (listings, handlers) can reuse it.
type ResolvePropertyTypesOutput struct {
	PropertyTypes    globalmodel.PropertyType
	ComplexName      string
	MainRegistration string
	Source           CoverageSource
}
