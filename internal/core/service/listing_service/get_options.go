package listingservices

import (
	"context"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ls *listingService) GetOptions(ctx context.Context, zipCode string, number string) (types []listingmodel.PropertyTypeOption, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	propertyTypes, err := ls.csi.GetOptions(ctx, zipCode, number)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	types = ls.DecodePropertyTypes(ctx, propertyTypes)

	return
}
