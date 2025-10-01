package globalmodel

// ContentSecurityPolicy represents the CSP directives configuration.
type ContentSecurityPolicy struct {
	ID         int64             `json:"id"`
	Version    int64             `json:"version"`
	Directives map[string]string `json:"directives"`
}

// NewContentSecurityPolicy creates a new immutable CSP entity.
func NewContentSecurityPolicy(id int64, version int64, directives map[string]string) ContentSecurityPolicy {
	normalized := make(map[string]string, len(directives))
	for key, value := range directives {
		normalized[key] = value
	}

	return ContentSecurityPolicy{
		ID:         id,
		Version:    version,
		Directives: normalized,
	}
}

// Clone returns a deep copy of the CSP policy, preserving immutability semantics.
func (c ContentSecurityPolicy) Clone() ContentSecurityPolicy {
	return NewContentSecurityPolicy(c.ID, c.Version, c.Directives)
}
