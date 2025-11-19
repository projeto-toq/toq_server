package listingservices

import (
	"context"
	"errors"
	"strings"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
	validators "github.com/projeto-toq/toq_server/internal/core/utils/validators"
)

func (ls *listingService) GetOptions(ctx context.Context, zipCode string, number string) (result listingmodel.PropertyOptionsResult, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return result, utils.InternalError("")
	}
	defer spanEnd()

	zipCandidate := strings.TrimSpace(zipCode)
	normalizedZip, normErr := validators.NormalizeCEP(zipCandidate)
	if normErr != nil {
		return result, utils.ValidationError("zipCode", "Zip code must contain exactly 8 digits without separators.")
	}
	zipCode = normalizedZip
	number = strings.TrimSpace(number)

	coverage, err := ls.propertyCoverage.ResolvePropertyTypes(ctx, propertycoveragemodel.ResolvePropertyTypesInput{
		ZipCode: zipCode,
		Number:  number,
	})
	if err != nil {
		utils.SetSpanError(ctx, err)
		var domainErr utils.DomainError
		if errors.As(err, &domainErr) {
			return result, utils.WrapDomainErrorWithSource(domainErr)
		}
		return result, utils.InternalError("")
	}

	result.PropertyTypes = ls.DecodePropertyTypes(ctx, coverage.PropertyTypes)
	result.ComplexName = coverage.ComplexName

	return result, nil
}
