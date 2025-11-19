package listingmodel

// PropertyOptionsResult aggregates the decoded property type options together
// with the complex name retrieved from property coverage lookups.
type PropertyOptionsResult struct {
	PropertyTypes []PropertyTypeOption
	ComplexName   string
}
