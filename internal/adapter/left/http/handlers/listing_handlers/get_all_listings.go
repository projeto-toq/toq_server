package listinghandlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetAllListings handles GET /listings
//
//	@Summary      List all listings with filters
//	@Description  Supports pagination, wildcard search for strings using '*', and range filters for numbers and dates.
//	@Tags         Listings
//	@Produce      json
//	@Param        page           query  int     false  "Page number" default(1) example(1)
//	@Param        limit          query  int     false  "Page size" default(10) example(20)
//	@Param        status         query  string  false  "Listing status (enum name or numeric)" example("PUBLISHED")
//	@Param        code           query  int     false  "Exact listing code" example(1024)
//	@Param        title          query  string  false  "Filter by listing title/description (supports '*' wildcard)" example("*garden*")
//	@Param        userId         query  int     false  "Filter by owner user id" example(55)
//	@Param        zipCode        query  string  false  "Filter by zip code (supports '*' wildcard)" example("12345*")
//	@Param        city           query  string  false  "Filter by city (supports '*' wildcard)" example("*Paulista*")
//	@Param        neighborhood   query  string  false  "Filter by neighborhood (supports '*' wildcard)" example("*Centro*")
//	@Param        createdFrom    query  string  false  "Filter by creation date from (RFC3339 or YYYY-MM-DD)" example("2025-01-01")
//	@Param        createdTo      query  string  false  "Filter by creation date to (RFC3339 or YYYY-MM-DD)" example("2025-01-31")
//	@Param        minSell        query  number  false  "Minimum sell price" example(100000)
//	@Param        maxSell        query  number  false  "Maximum sell price" example(900000)
//	@Param        minRent        query  number  false  "Minimum rent price" example(1500)
//	@Param        maxRent        query  number  false  "Maximum rent price" example(8000)
//	@Param        minLandSize    query  number  false  "Minimum land size" example(120.5)
//	@Param        maxLandSize    query  number  false  "Maximum land size" example(500.75)
//	@Success      200  {object}  dto.GetAllListingsResponse
//	@Failure      400  {object}  map[string]any
//	@Failure      401  {object}  map[string]any
//	@Failure      403  {object}  map[string]any
//	@Failure      500  {object}  map[string]any
//	@Router       /listings [get]
func (lh *ListingHandler) GetAllListings(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	var req dto.GetAllListingsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httperrors.SendHTTPErrorObj(c, httperrors.ConvertBindError(err))
		return
	}

	statusPtr, err := parseListingStatus(strings.TrimSpace(req.Status))
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("status", err.Error()))
		return
	}

	codePtr, err := parseOptionalUint32(req.Code)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("code", err.Error()))
		return
	}

	userIDPtr, err := parseOptionalInt64(req.UserID)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("userId", err.Error()))
		return
	}

	createdFrom, err := parseOptionalTime(req.CreatedFrom)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("createdFrom", err.Error()))
		return
	}
	createdTo, err := parseOptionalTime(req.CreatedTo)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("createdTo", err.Error()))
		return
	}
	if createdFrom != nil && createdTo != nil && createdFrom.After(*createdTo) {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("createdFrom", "from date cannot be greater than to date"))
		return
	}

	minSell, err := parseOptionalFloat64(req.MinSellPrice)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("minSell", err.Error()))
		return
	}
	maxSell, err := parseOptionalFloat64(req.MaxSellPrice)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("maxSell", err.Error()))
		return
	}
	if minSell != nil && maxSell != nil && *minSell > *maxSell {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("minSell", "minSell cannot be greater than maxSell"))
		return
	}

	minRent, err := parseOptionalFloat64(req.MinRentPrice)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("minRent", err.Error()))
		return
	}
	maxRent, err := parseOptionalFloat64(req.MaxRentPrice)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("maxRent", err.Error()))
		return
	}
	if minRent != nil && maxRent != nil && *minRent > *maxRent {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("minRent", "minRent cannot be greater than maxRent"))
		return
	}

	minLand, err := parseOptionalFloat64(req.MinLandSize)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("minLandSize", err.Error()))
		return
	}
	maxLand, err := parseOptionalFloat64(req.MaxLandSize)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("maxLandSize", err.Error()))
		return
	}
	if minLand != nil && maxLand != nil && *minLand > *maxLand {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("minLandSize", "minLandSize cannot be greater than maxLandSize"))
		return
	}

	input := listingservices.ListListingsInput{
		Page:         req.Page,
		Limit:        req.Limit,
		Status:       statusPtr,
		Code:         codePtr,
		Title:        strings.TrimSpace(req.Title),
		ZipCode:      strings.TrimSpace(req.ZipCode),
		City:         strings.TrimSpace(req.City),
		Neighborhood: strings.TrimSpace(req.Neighborhood),
		UserID:       userIDPtr,
		CreatedFrom:  createdFrom,
		CreatedTo:    createdTo,
		MinSellPrice: minSell,
		MaxSellPrice: maxSell,
		MinRentPrice: minRent,
		MaxRentPrice: maxRent,
		MinLandSize:  minLand,
		MaxLandSize:  maxLand,
	}

	result, listErr := lh.listingService.ListListings(ctx, input)
	if listErr != nil {
		httperrors.SendHTTPErrorObj(c, listErr)
		return
	}

	data := make([]dto.ListingResponse, 0, len(result.Items))
	for _, item := range result.Items {
		data = append(data, toListingResponse(item))
	}

	resp := dto.GetAllListingsResponse{
		Data: data,
		Pagination: dto.PaginationResponse{
			Page:       result.Page,
			Limit:      result.Limit,
			Total:      result.Total,
			TotalPages: computeTotalPages(result.Total, result.Limit),
		},
	}

	c.JSON(http.StatusOK, resp)
}

func parseListingStatus(raw string) (*listingmodel.ListingStatus, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}

	if numeric, numErr := strconv.Atoi(trimmed); numErr == nil {
		status := listingmodel.ListingStatus(numeric)
		if status >= listingmodel.StatusDraft && status <= listingmodel.StatusArchived {
			return &status, nil
		}
		return nil, fmt.Errorf("invalid listing status value")
	}

	upper := strings.ToUpper(trimmed)
	normalized := strings.ReplaceAll(strings.ReplaceAll(upper, " ", ""), "-", "")
	normalized = strings.ReplaceAll(normalized, "_", "")

	mapping := map[string]listingmodel.ListingStatus{
		"DRAFT":                  listingmodel.StatusDraft,
		"PENDINGPHOTOSCHEDULING": listingmodel.StatusPendingPhotoScheduling,
		"PHOTOSSCHEDULED":        listingmodel.StatusPhotosScheduled,
		"PENDINGPHOTOPROCESSING": listingmodel.StatusPendingPhotoProcessing,
		"PENDINGOWNERAPPROVAL":   listingmodel.StatusPendingOwnerApproval,
		"PENDINGADMINREVIEW":     listingmodel.StatusPendingAdminReview,
		"PUBLISHED":              listingmodel.StatusPublished,
		"UNDEROFFER":             listingmodel.StatusUnderOffer,
		"UNDERNEGOTIATION":       listingmodel.StatusUnderNegotiation,
		"CLOSED":                 listingmodel.StatusClosed,
		"SUSPENDED":              listingmodel.StatusSuspended,
		"REJECTEDBYOWNER":        listingmodel.StatusRejectedByOwner,
		"NEEDSREVISION":          listingmodel.StatusNeedsRevision,
		"EXPIRED":                listingmodel.StatusExpired,
		"ARCHIVED":               listingmodel.StatusArchived,
	}
	if status, ok := mapping[normalized]; ok {
		return &status, nil
	}

	return nil, fmt.Errorf("invalid listing status")
}

func parseOptionalUint32(raw string) (*uint32, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}
	value, err := strconv.ParseUint(trimmed, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid numeric value")
	}
	v := uint32(value)
	return &v, nil
}

func parseOptionalInt64(raw string) (*int64, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}
	value, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid numeric value")
	}
	return &value, nil
}

func parseOptionalFloat64(raw string) (*float64, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}
	value, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid decimal value")
	}
	return &value, nil
}

func parseOptionalTime(raw string) (*time.Time, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}

	layouts := []string{time.RFC3339Nano, time.RFC3339, "2006-01-02"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, trimmed); err == nil {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("invalid datetime format")
}

func toListingResponse(item listingservices.ListListingsItem) dto.ListingResponse {
	listing := item.Listing

	price := listing.SellNet()
	if price == 0 {
		price = listing.RentNet()
	}

	return dto.ListingResponse{
		ID:           listing.ID(),
		Title:        deriveListingTitle(listing),
		Description:  strings.TrimSpace(listing.Description()),
		Price:        price,
		Status:       listing.Status().String(),
		PropertyType: int(listing.ListingType()),
		ZipCode:      listing.ZipCode(),
		Number:       listing.Number(),
		UserID:       listing.UserID(),
		CreatedAt:    formatTime(item.CreatedAt),
		UpdatedAt:    formatTime(item.UpdatedAt),
	}
}

func deriveListingTitle(listing listingmodel.ListingInterface) string {
	description := strings.TrimSpace(listing.Description())
	if description != "" {
		if len(description) > 120 {
			description = description[:120]
		}
		return description
	}

	if listing.Code() != 0 {
		return fmt.Sprintf("Listing %d", listing.Code())
	}

	if listing.ID() != 0 {
		return fmt.Sprintf("Listing %d", listing.ID())
	}

	return "Listing"
}

func formatTime(ts *time.Time) string {
	if ts == nil {
		return ""
	}
	return ts.UTC().Format(time.RFC3339)
}

func computeTotalPages(total int64, limit int) int {
	if limit <= 0 || total <= 0 {
		return 0
	}

	pages := int(total / int64(limit))
	if total%int64(limit) != 0 {
		pages++
	}

	if pages == 0 && total > 0 {
		return 1
	}

	return pages
}
