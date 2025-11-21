package listingservices

import (
	"context"
	"strings"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	propertycoverageservice "github.com/projeto-toq/toq_server/internal/core/service/property_coverage_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
	validators "github.com/projeto-toq/toq_server/internal/core/utils/validators"
)

// GetOptions retrieves the full complex details or coverage options for a given address.
// It delegates to PropertyCoverageService.GetComplexByAddress to handle both complex and standalone scenarios.
func (ls *listingService) GetOptions(ctx context.Context, zipCode string, number string) (propertycoveragemodel.ManagedComplexInterface, error) {
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

	// Delegate to PropertyCoverageService which now handles Vertical, Horizontal and Standalone cases
	complex, err := ls.propertyCoverage.GetComplexByAddress(ctx, propertycoverageservice.GetComplexByAddressInput{
		ZipCode: zipCode,
		Number:  number,
	})
	if err != nil {
		utils.SetSpanError(ctx, err)
		return nil, err
	}

	return complex, nil
}
