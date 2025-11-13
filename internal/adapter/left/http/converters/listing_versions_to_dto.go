package converters

import (
	"strings"

	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
)

// ListingVersionsToDTO converte a saída do service para o DTO exposto em HTTP.
func ListingVersionsToDTO(output listingservices.ListListingVersionsOutput) dto.ListListingVersionsResponse {
	response := dto.ListListingVersionsResponse{
		Versions: make([]dto.ListingVersionSummaryResponse, 0, len(output.Versions)),
	}

	for _, info := range output.Versions {
		version := info.Version
		if version == nil {
			continue
		}

		if response.ListingIdentityID == 0 {
			response.ListingIdentityID = version.ListingIdentityID()
		}
		if response.ListingUUID == "" {
			response.ListingUUID = version.ListingUUID()
		}

		title := strings.TrimSpace(version.Title())

		response.Versions = append(response.Versions, dto.ListingVersionSummaryResponse{
			ID:                version.ID(),
			ListingIdentityID: version.ListingIdentityID(),
			ListingUUID:       version.ListingUUID(),
			Version:           version.Version(),
			Status:            version.Status().String(),
			Title:             title,
			IsActive:          info.IsActive,
		})
	}

	return response
}
