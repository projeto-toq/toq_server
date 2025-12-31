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

	visitType, err := parseVisitType(req.Type)
	if err != nil {
		return visitservice.CreateVisitInput{}, err
	}

	input := visitservice.CreateVisitInput{
		ListingIdentityID: req.ListingIdentityID,
		ScheduledStart:    start,
		ScheduledEnd:      end,
		Type:              visitType,
		RealtorNotes:      strings.TrimSpace(req.RealtorNotes),
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
		DurationMinutes:   visit.DurationMinutes(),
		Status:            string(visit.Status()),
		Type:              string(visit.Type()),
		CreatedAt:         visit.CreatedAt().Format(time.RFC3339),
		UpdatedAt:         visit.UpdatedAt().Format(time.RFC3339),
	}

	if source, ok := visit.Source(); ok {
		response.Source = source
	}

	if notes, ok := visit.RealtorNotes(); ok {
		response.RealtorNotes = notes
	}

	if notes, ok := visit.OwnerNotes(); ok {
		response.OwnerNotes = notes
	}

	if reason, ok := visit.RejectionReason(); ok {
		response.RejectionReason = reason
	}

	if reason, ok := visit.CancelReason(); ok {
		response.CancelReason = reason
	}

	if ts, ok := visit.FirstOwnerActionAt(); ok {
		formatted := ts.Format(time.RFC3339)
		response.FirstOwnerActionAt = &formatted
	}

	return response
}

// VisitListToResponse builds a paginated response DTO.
func VisitListToResponse(result listingmodel.VisitListResult, page, limit int) dto.VisitListResponse {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	items := make([]dto.VisitResponse, 0, len(result.Visits))
	for _, visit := range result.Visits {
		items = append(items, VisitDomainToResponse(visit))
	}

	return dto.VisitListResponse{
		Items: items,
		Pagination: dto.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      result.Total,
			TotalPages: visitTotalPages(result.Total, limit),
		},
	}
}

func parseVisitType(raw string) (listingmodel.VisitMode, error) {
	normalized := strings.ToUpper(strings.TrimSpace(raw))
	visitType := listingmodel.VisitMode(normalized)
	switch visitType {
	case listingmodel.VisitModeWithClient, listingmodel.VisitModeRealtorOnly, listingmodel.VisitModeContentProduction:
		return visitType, nil
	default:
		return "", coreutils.ValidationError("type", "invalid visit type")
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
