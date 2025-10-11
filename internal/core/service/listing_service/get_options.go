package listingservices

import (
	"context"
	"errors"

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
		var domainErr utils.DomainError
		if errors.As(err, &domainErr) {
			return nil, utils.WrapDomainErrorWithSource(domainErr)
		}
		return nil, utils.InternalError("")
	}

	types = ls.DecodePropertyTypes(ctx, propertyTypes)

	return
}
