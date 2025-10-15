package listingservices

import (
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

type CreateCatalogValueInput struct {
	Category    string
	Slug        string
	Label       string
	Description coreutils.Optional[string]
	IsActive    coreutils.Optional[bool]
}

type UpdateCatalogValueInput struct {
	Category    string
	ID          uint8
	Slug        coreutils.Optional[string]
	Label       coreutils.Optional[string]
	Description coreutils.Optional[string]
	IsActive    coreutils.Optional[bool]
}

type RestoreCatalogValueInput struct {
	Category string
	ID       uint8
}

func (in CreateCatalogValueInput) ToDomain(nextID uint8) listingmodel.CatalogValueInterface {
	value := listingmodel.NewCatalogValue()
	value.SetID(nextID)
	value.SetCategory(in.Category)
	value.SetSlug(in.Slug)
	value.SetLabel(in.Label)
	if in.Description.IsPresent() {
		if in.Description.IsNull() {
			value.SetDescription(nil)
		} else if desc, ok := in.Description.Value(); ok {
			copyDesc := desc
			value.SetDescription(&copyDesc)
		}
	}
	active := true
	if in.IsActive.IsPresent() {
		if in.IsActive.IsNull() {
			active = false
		} else if val, ok := in.IsActive.Value(); ok {
			active = val
		}
	}
	value.SetIsActive(active)
	return value
}

func (in UpdateCatalogValueInput) ApplyToDomain(value listingmodel.CatalogValueInterface) {
	if in.Slug.IsPresent() {
		if in.Slug.IsNull() {
			value.SetSlug("")
		} else if slug, ok := in.Slug.Value(); ok {
			value.SetSlug(slug)
		}
	}

	if in.Label.IsPresent() {
		if in.Label.IsNull() {
			value.SetLabel("")
		} else if label, ok := in.Label.Value(); ok {
			value.SetLabel(label)
		}
	}

	if in.Description.IsPresent() {
		if in.Description.IsNull() {
			value.SetDescription(nil)
		} else if desc, ok := in.Description.Value(); ok {
			copyDesc := desc
			value.SetDescription(&copyDesc)
		}
	}

	if in.IsActive.IsPresent() {
		if in.IsActive.IsNull() {
			value.SetIsActive(false)
		} else if active, ok := in.IsActive.Value(); ok {
			value.SetIsActive(active)
		}
	}
}
