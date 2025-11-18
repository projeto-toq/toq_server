package listingmodel

import "strconv"

// ListingIdentityID representa o identificador lógico (listing_identities.id).
type ListingIdentityID uint64

// ListingVersionID representa uma versão específica (listing_versions.id).
type ListingVersionID uint64

// Uint64 retorna o valor bruto para uso em adapters/infra.
func (id ListingIdentityID) Uint64() uint64 { return uint64(id) }

// Int64 retorna o valor como int64 (útil para queries legacy).
func (id ListingIdentityID) Int64() int64 { return int64(id) }

// String retorna o ID em formato decimal.
func (id ListingIdentityID) String() string { return strconv.FormatUint(uint64(id), 10) }

// IsZero indica se o ID foi preenchido.
func (id ListingIdentityID) IsZero() bool { return id == 0 }

// ListingIdentityIDFromInt64 converte valores assinados de forma segura.
func ListingIdentityIDFromInt64(value int64) ListingIdentityID {
	if value <= 0 {
		return 0
	}
	return ListingIdentityID(uint64(value))
}

// Uint64 expõe o valor bruto do version ID.
func (id ListingVersionID) Uint64() uint64 { return uint64(id) }

// Int64 expõe o valor do version ID como int64.
func (id ListingVersionID) Int64() int64 { return int64(id) }

// String retorna a representação decimal do version ID.
func (id ListingVersionID) String() string { return strconv.FormatUint(uint64(id), 10) }

// IsZero indica se o version ID foi definido.
func (id ListingVersionID) IsZero() bool { return id == 0 }

// ListingVersionIDFromInt64 converte valores assinados para version IDs fortes.
func ListingVersionIDFromInt64(value int64) ListingVersionID {
	if value <= 0 {
		return 0
	}
	return ListingVersionID(uint64(value))
}
