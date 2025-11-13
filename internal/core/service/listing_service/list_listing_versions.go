package listingservices

import (
	"context"

	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ls *listingService) ListListingVersions(ctx context.Context, input ListListingVersionsInput) (ListListingVersionsOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return ListListingVersionsOutput{}, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.ListingIdentityID <= 0 {
		return ListListingVersionsOutput{}, utils.ValidationError("listingIdentityId", "Listing identity id must be greater than zero")
	}

	tx, txErr := ls.gsi.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("listing.list_versions.tx_start_error", "error", txErr, "identity_id", input.ListingIdentityID)
		return ListListingVersionsOutput{}, utils.InternalError("")
	}
	defer func() {
		_ = ls.gsi.RollbackTransaction(ctx, tx)
	}()

	filter := listingrepository.ListListingVersionsFilter{
		ListingIdentityID: input.ListingIdentityID,
		IncludeDeleted:    input.IncludeDeleted,
	}

	summaries, listErr := ls.listingRepository.ListListingVersions(ctx, tx, filter)
	if listErr != nil {
		utils.SetSpanError(ctx, listErr)
		logger.Error("listing.list_versions.repo_error", "error", listErr, "identity_id", input.ListingIdentityID)
		return ListListingVersionsOutput{}, utils.InternalError("")
	}

	output := ListListingVersionsOutput{Versions: make([]ListingVersionInfo, 0, len(summaries))}
	for _, summary := range summaries {
		output.Versions = append(output.Versions, ListingVersionInfo{
			Version:  summary.Version,
			IsActive: summary.IsActive,
		})
	}

	return output, nil
}
