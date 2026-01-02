package visithandlers

import (
	"strings"
	"time"

	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

const maxVisitPageSize = 50

func buildVisitFilter(query dto.VisitListQuery, actorID int64, asOwner bool) (listingmodel.VisitListFilter, error) {
	statuses, err := parseVisitStatuses(query.Statuses)
	if err != nil {
		return listingmodel.VisitListFilter{}, err
	}

	from, err := parseOptionalTime("from", query.From)
	if err != nil {
		return listingmodel.VisitListFilter{}, err
	}

	to, err := parseOptionalTime("to", query.To)
	if err != nil {
		return listingmodel.VisitListFilter{}, err
	}

	if from != nil && to != nil && from.After(*to) {
		return listingmodel.VisitListFilter{}, coreutils.ValidationError("from", "must be before or equal to 'to'")
	}

	page := query.Page
	limit := query.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > maxVisitPageSize {
		limit = maxVisitPageSize
	}

	filter := listingmodel.VisitListFilter{
		Statuses: statuses,
		From:     from,
		To:       to,
		Page:     page,
		Limit:    limit,
	}

	if query.ListingIdentityID > 0 {
		filter.ListingIdentityID = &query.ListingIdentityID
	}

	if asOwner {
		filter.OwnerUserID = &actorID
	} else {
		filter.RequesterUserID = &actorID
	}

	return filter, nil
}

func parseVisitStatuses(raw []string) ([]listingmodel.VisitStatus, error) {
	normalized := normalizeMulti(raw)
	statuses := make([]listingmodel.VisitStatus, 0, len(normalized))
	for _, item := range normalized {
		status, err := listingmodel.ParseVisitStatus(item)
		if err != nil {
			return nil, coreutils.ValidationError("status", err.Error())
		}
		statuses = append(statuses, status)
	}
	return statuses, nil
}

func parseOptionalTime(field, value string) (*time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}
	ts, err := time.Parse(time.RFC3339, trimmed)
	if err != nil {
		return nil, coreutils.ValidationError(field, "must be a valid RFC3339 timestamp")
	}
	return &ts, nil
}

func normalizeMulti(raw []string) []string {
	result := make([]string, 0)
	for _, item := range raw {
		for _, part := range strings.Split(item, ",") {
			clean := strings.TrimSpace(part)
			if clean == "" {
				continue
			}
			result = append(result, clean)
		}
	}
	return result
}
