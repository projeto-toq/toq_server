package listingservices

import (
	"context"
	"strings"
	"time"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListListingsInput captures filters and pagination for listing search.
type ListListingsInput struct {
	Page              int
	Limit             int
	Status            *listingmodel.ListingStatus
	Code              *uint32
	Title             string
	ZipCode           string
	City              string
	Neighborhood      string
	UserID            *int64
	CreatedFrom       *time.Time
	CreatedTo         *time.Time
	MinSellPrice      *float64
	MaxSellPrice      *float64
	MinRentPrice      *float64
	MaxRentPrice      *float64
	MinLandSize       *float64
	MaxLandSize       *float64
	RequesterUserID   int64
	RequesterRoleSlug permissionmodel.RoleSlug
}

// ListListingsOutput encapsulates listings and paging metadata.
type ListListingsOutput struct {
	Items []ListListingsItem
	Total int64
	Page  int
	Limit int
}

// ListListingsItem traz a entidade de listing e metadados úteis para montagem da resposta.
type ListListingsItem struct {
	Listing   listingmodel.ListingInterface
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

// ListListings returns listings filtered with pagination for admin panel consumption.
func (ls *listingService) ListListings(ctx context.Context, input ListListingsInput) (ListListingsOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return ListListingsOutput{}, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.RequesterRoleSlug == permissionmodel.RoleSlugOwner {
		ownerID := input.RequesterUserID
		input.UserID = &ownerID
		logger.Debug("listing.list.owner_scope_enforced", "user_id", ownerID)
	}

	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 20
	}

	tx, txErr := ls.gsi.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("listing.list.tx_start_failed", "error", txErr)
		return ListListingsOutput{}, utils.InternalError("")
	}
	defer func() {
		_ = ls.gsi.RollbackTransaction(ctx, tx)
	}()

	zipFilter := sanitizeZipFilter(input.ZipCode)

	repoFilter := listingrepository.ListListingsFilter{
		Page:         input.Page,
		Limit:        input.Limit,
		Status:       input.Status,
		Code:         input.Code,
		Title:        utils.NormalizeSearchPattern(input.Title),
		ZipCode:      utils.NormalizeSearchPattern(zipFilter),
		City:         utils.NormalizeSearchPattern(input.City),
		Neighborhood: utils.NormalizeSearchPattern(input.Neighborhood),
		UserID:       input.UserID,
		CreatedFrom:  input.CreatedFrom,
		CreatedTo:    input.CreatedTo,
		MinSellPrice: input.MinSellPrice,
		MaxSellPrice: input.MaxSellPrice,
		MinRentPrice: input.MinRentPrice,
		MaxRentPrice: input.MaxRentPrice,
		MinLandSize:  input.MinLandSize,
		MaxLandSize:  input.MaxLandSize,
	}

	result, listErr := ls.listingRepository.ListListings(ctx, tx, repoFilter)
	if listErr != nil {
		utils.SetSpanError(ctx, listErr)
		logger.Error("listing.list.repo_error", "error", listErr)
		return ListListingsOutput{}, utils.InternalError("")
	}

	items := make([]ListListingsItem, 0, len(result.Records))
	for _, record := range result.Records {
		items = append(items, ListListingsItem{
			Listing:   record.Listing,
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt,
		})
	}

	return ListListingsOutput{
		Items: items,
		Total: result.Total,
		Page:  repoFilter.Page,
		Limit: repoFilter.Limit,
	}, nil
}

func sanitizeZipFilter(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}
	var builder strings.Builder
	for _, r := range trimmed {
		switch {
		case r >= '0' && r <= '9':
			builder.WriteRune(r)
		case r == '*', r == '%':
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
