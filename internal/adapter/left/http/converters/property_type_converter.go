package converters

import (
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

// BuildListingPropertyTypeDTO converts a domain property type bitmask into the enriched response DTO.
func BuildListingPropertyTypeDTO(propertyType globalmodel.PropertyType) *dto.ListingPropertyTypeResponse {
	option, ok := listingmodel.PropertyTypeOptionFromBit(propertyType)
	if !ok {
		return nil
	}

	return &dto.ListingPropertyTypeResponse{
		Code:        option.Code,
		Label:       option.Label,
		PropertyBit: uint16(option.PropertyBit),
	}
}
