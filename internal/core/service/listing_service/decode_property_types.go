package listingservices

import (
	"context"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
)

func (ls *listingService) DecodePropertyTypes(ctx context.Context, propertyTypes globalmodel.PropertyType) (types []int64) {

	if propertyTypes&globalmodel.Apartment == globalmodel.Apartment {
		types = append(types, 1)
	}

	if propertyTypes&globalmodel.CommercialStore == globalmodel.CommercialStore {
		types = append(types, 2)
	}

	if propertyTypes&globalmodel.CommercialFloor == globalmodel.CommercialFloor {
		types = append(types, 4)
	}

	if propertyTypes&globalmodel.Suite == globalmodel.Suite {
		types = append(types, 8)
	}

	if propertyTypes&globalmodel.House == globalmodel.House {
		types = append(types, 16)
	}

	if propertyTypes&globalmodel.OffPlanHouse == globalmodel.OffPlanHouse {
		types = append(types, 32)
	}

	if propertyTypes&globalmodel.ResidencialLand == globalmodel.ResidencialLand {
		types = append(types, 64)
	}

	if propertyTypes&globalmodel.CommercialLand == globalmodel.CommercialLand {
		types = append(types, 128)
	}

	if propertyTypes&globalmodel.Building == globalmodel.Building {
		types = append(types, 256)
	}

	if propertyTypes&globalmodel.Warehouse == globalmodel.Warehouse {
		types = append(types, 512)
	}

	return
}
