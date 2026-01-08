package converters

import (
	"strings"
	"time"

	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	visitservice "github.com/projeto-toq/toq_server/internal/core/service/visit_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateVisitDTOToInput converts the incoming DTO into a service input with parsed values.
func CreateVisitDTOToInput(req dto.CreateVisitRequest) (visitservice.CreateVisitInput, error) {
	start, err := time.Parse(time.RFC3339, strings.TrimSpace(req.ScheduledStart))
	if err != nil {
		return visitservice.CreateVisitInput{}, coreutils.ValidationError("scheduledStart", "must be a valid RFC3339 timestamp")
	}

	end, err := time.Parse(time.RFC3339, strings.TrimSpace(req.ScheduledEnd))
	if err != nil {
		return visitservice.CreateVisitInput{}, coreutils.ValidationError("scheduledEnd", "must be a valid RFC3339 timestamp")
	}

	if !start.Before(end) {
		return visitservice.CreateVisitInput{}, coreutils.ValidationError("scheduledStart", "must be before scheduledEnd")
	}

	input := visitservice.CreateVisitInput{
		ListingIdentityID: req.ListingIdentityID,
		ScheduledStart:    start,
		ScheduledEnd:      end,
		Notes:             strings.TrimSpace(req.Notes),
		Source:            strings.TrimSpace(req.Source),
	}

	return input, nil
}

// VisitDomainToResponse maps the domain model to a response DTO.
func VisitDomainToResponse(visit listingmodel.VisitInterface) dto.VisitResponse {
	if visit == nil {
		return dto.VisitResponse{}
	}

	response := dto.VisitResponse{
		ID:                visit.ID(),
		ListingIdentityID: visit.ListingIdentityID(),
		ListingVersion:    visit.ListingVersion(),
		RequesterUserID:   visit.RequesterUserID(),
		OwnerUserID:       visit.OwnerUserID(),
		ScheduledStart:    visit.ScheduledStart().Format(time.RFC3339),
		ScheduledEnd:      visit.ScheduledEnd().Format(time.RFC3339),
		Status:            string(visit.Status()),
	}

	if source, ok := visit.Source(); ok {
		response.Source = source
	}

	if notes, ok := visit.Notes(); ok {
		response.Notes = notes
	}

	if reason, ok := visit.RejectionReason(); ok {
		response.RejectionReason = reason
	}

	if ts, ok := visit.FirstOwnerActionAt(); ok {
		formatted := ts.Format(time.RFC3339)
		response.FirstOwnerActionAt = &formatted
	}

	return response
}

// VisitDetailToResponse enriches the visit response with the related listing snapshot.
func VisitDetailToResponse(detail visitservice.VisitDetailOutput) dto.VisitResponse {
	response := VisitDomainToResponse(detail.Visit)
	listing := detail.Listing
	if listing == nil {
		return response
	}

	summary := dto.ListingSummaryDTO{
		Title:        strings.TrimSpace(listing.Title()),
		Description:  strings.TrimSpace(listing.Description()),
		ZipCode:      listing.ZipCode(),
		Street:       listing.Street(),
		Number:       listing.Number(),
		Neighborhood: listing.Neighborhood(),
		City:         listing.City(),
		State:        listing.State(),
	}

	if complement := strings.TrimSpace(listing.Complement()); complement != "" {
		summary.Complement = complement
	}

	response.ListingSummary = &summary
	return response
}

// VisitListDetailToResponse builds a paginated response DTO using hydrated listing snapshots.
func VisitListDetailToResponse(output visitservice.VisitListOutput) dto.VisitListResponse {
	page := output.Page
	limit := output.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	items := make([]dto.VisitResponse, 0, len(output.Items))
	for _, detail := range output.Items {
		items = append(items, VisitDetailToResponse(detail))
	}

	return dto.VisitListResponse{
		Items: items,
		Pagination: dto.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      output.Total,
			TotalPages: visitTotalPages(output.Total, limit),
		},
	}
}

func visitTotalPages(total int64, limit int) int {
	if limit <= 0 {
		return 0
	}
	if total == 0 {
		return 0
	}

	pages := total / int64(limit)
	if total%int64(limit) != 0 {
		pages++
	}

	return int(pages)
}
