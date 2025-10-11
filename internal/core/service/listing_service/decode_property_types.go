package listingservices

import (
	"context"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

func (ls *listingService) DecodePropertyTypes(ctx context.Context, propertyTypes globalmodel.PropertyType) (types []listingmodel.PropertyTypeOption) {
	for _, option := range listingmodel.PropertyTypeCatalog {
		if propertyTypes&option.PropertyBit == option.PropertyBit {
			types = append(types, option)
		}
	}

	return
}
