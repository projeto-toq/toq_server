package listingservices

import (
	"context"
	"errors"
	"strings"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
	validators "github.com/projeto-toq/toq_server/internal/core/utils/validators"
)

func (ls *listingService) GetOptions(ctx context.Context, zipCode string, number string) (types []listingmodel.PropertyTypeOption, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	zipCandidate := strings.TrimSpace(zipCode)
	normalizedZip, normErr := validators.NormalizeCEP(zipCandidate)
	if normErr != nil {
		return nil, utils.ValidationError("zipCode", "Zip code must contain exactly 8 digits without separators.")
	}
	zipCode = normalizedZip
	number = strings.TrimSpace(number)

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
