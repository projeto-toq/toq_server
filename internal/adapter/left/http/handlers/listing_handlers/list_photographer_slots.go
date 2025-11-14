package listinghandlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

const (
	isoDateLayout        = "2006-01-02"
	fallbackSlotsPerPage = 20
)

// ListPhotographerSlots lista slots disponíveis na agenda de fotógrafos.
//
//	@Summary   List available photographer slots
//	@Tags      Listing Photo Sessions
//	@Accept    json
//	@Produce   json
//	@Param     from      query    string false "Start date filter (YYYY-MM-DD)" Format(date) example(2025-10-20)
//	@Param     to        query    string false "End date filter (YYYY-MM-DD)" Format(date) example(2025-10-31)
//	@Param     period    query    string false "Slot period" Enums(MORNING,AFTERNOON) example(MORNING)
//	@Param     page      query    int    false "Page number" default(1)
//	@Param     size      query    int    false "Page size" default(20)
//	@Param     sort      query    string false "Sort order" Enums(start_asc,start_desc,photographer_asc,photographer_desc) default(start_asc)
//	@Param     listingIdentityId query    int    true  "Listing identifier" example(1024)
//	@Param     timezone  query    string true  "Listing timezone" example(America/Sao_Paulo)
//	@Success   200 {object} dto.ListPhotographerSlotsResponse
//	@Failure   400 {object} dto.ErrorResponse "Invalid filters"
//	@Failure   401 {object} dto.ErrorResponse "Unauthorized"
//	@Failure   403 {object} dto.ErrorResponse "Forbidden"
//	@Failure   500 {object} dto.ErrorResponse "Internal error"
//	@Router    /listings/photo-session/slots [get]
//	@Security  BearerAuth
func (lh *ListingHandler) ListPhotographerSlots(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	ctx, spanEnd, err := coreutils.GenerateTracer(baseCtx)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	if _, ok := middlewares.GetUserInfoFromContext(c); !ok {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_CONTEXT_MISSING", "User context not found")
		return
	}

	var request dto.ListPhotographerSlotsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid filter parameters")
		return
	}

	timezone := strings.TrimSpace(request.Timezone)
	if timezone == "" {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_TIMEZONE", "timezone is required")
		return
	}

	loc, tzErr := time.LoadLocation(timezone)
	if tzErr != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_TIMEZONE", "timezone must be a valid IANA identifier")
		return
	}

	var fromPtr *time.Time
	if request.From != "" {
		parsed, parseErr := time.Parse(isoDateLayout, request.From)
		if parseErr != nil {
			httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_FROM_DATE", "from must be in YYYY-MM-DD format")
			return
		}
		local := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, loc)
		fromPtr = &local
	}

	var toPtr *time.Time
	if request.To != "" {
		parsed, parseErr := time.Parse(isoDateLayout, request.To)
		if parseErr != nil {
			httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_TO_DATE", "to must be in YYYY-MM-DD format")
			return
		}
		local := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, loc)
		toPtr = &local
	}

	var periodPtr *photosessionmodel.SlotPeriod
	if request.Period != "" {
		period := photosessionmodel.SlotPeriod(request.Period)
		periodPtr = &period
	}

	input := listingservices.ListPhotographerSlotsInput{
		From:              fromPtr,
		To:                toPtr,
		Period:            periodPtr,
		Page:              request.Page,
		Size:              request.Size,
		Sort:              strings.TrimSpace(request.Sort),
		ListingIdentityID: request.ListingIdentityID,
		Location:          loc,
	}

	output, err := lh.listingService.ListPhotographerSlots(ctx, input)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	slots := make([]dto.PhotographerSlotResponse, 0, len(output.Slots))
	for _, slot := range output.Slots {
		var reservedUntil *string
		if until := slot.ReservedUntil(); until != nil {
			formatted := until.UTC().Format(time.RFC3339)
			reservedUntil = &formatted
		}

		slots = append(slots, dto.PhotographerSlotResponse{
			SlotID:             slot.ID(),
			PhotographerUserID: slot.PhotographerUserID(),
			SlotStart:          slot.SlotStart().UTC().Format(time.RFC3339),
			SlotEnd:            slot.SlotEnd().UTC().Format(time.RFC3339),
			Status:             string(slot.Status()),
			ReservedUntil:      reservedUntil,
		})
	}

	size := output.Size
	if size <= 0 {
		size = fallbackSlotsPerPage
	}

	totalPages := 0
	if output.Total > 0 {
		totalPages = int((output.Total + int64(size) - 1) / int64(size))
	}

	response := dto.ListPhotographerSlotsResponse{
		Data: slots,
		Pagination: dto.PaginationResponse{
			Page:       output.Page,
			Limit:      size,
			Total:      output.Total,
			TotalPages: totalPages,
		},
	}

	c.JSON(http.StatusOK, response)
}
