package photosessionservices

import (
	"context"
	"sort"
	"strings"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListServiceAreas returns paginated service areas filtered by city and state.
func (s *photoSessionService) ListServiceAreas(ctx context.Context, input ListServiceAreasInput) (ListServiceAreasOutput, error) {
	if input.PhotographerID == 0 {
		return ListServiceAreasOutput{}, utils.ValidationError("photographerId", "photographerId must be greater than zero")
	}

	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "service.ListServiceAreas")
	if err != nil {
		return ListServiceAreasOutput{}, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	page := input.Page
	if page <= 0 {
		page = defaultServiceAreaPage
	}

	size := input.Size
	if size <= 0 {
		size = defaultServiceAreaSize
	}
	if size > maxServiceAreaPageSize {
		size = maxServiceAreaPageSize
	}

	tx, err := s.globalService.StartReadOnlyTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.service_area.list.tx_start_error", "err", err)
		return ListServiceAreasOutput{}, utils.InternalError("")
	}
	defer func() {
		if rollbackErr := s.globalService.RollbackTransaction(ctx, tx); rollbackErr != nil {
			utils.SetSpanError(ctx, rollbackErr)
			logger.Error("photo_session.service_area.list.tx_rollback_error", "err", rollbackErr)
		}
	}()

	areas, err := s.repo.ListServiceAreasByPhotographer(ctx, tx, input.PhotographerID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.service_area.list.repo_error", "photographer_id", input.PhotographerID, "err", err)
		return ListServiceAreasOutput{}, utils.InternalError("")
	}

	filtered := filterServiceAreas(areas, input.City, input.State)
	sortServiceAreas(filtered)

	total := len(filtered)
	start := (page - 1) * size
	if start > total {
		start = total
	}
	end := start + size
	if end > total {
		end = total
	}

	return ListServiceAreasOutput{
		Areas: filtered[start:end],
		Total: int64(total),
		Page:  page,
		Size:  size,
	}, nil
}

func filterServiceAreas(source []photosessionmodel.PhotographerServiceAreaInterface, cityFilter, stateFilter *string) []photosessionmodel.PhotographerServiceAreaInterface {
	if len(source) == 0 {
		return source
	}

	normalizedCity := ""
	if cityFilter != nil {
		normalizedCity = strings.TrimSpace(*cityFilter)
	}

	normalizedState := ""
	if stateFilter != nil {
		normalizedState = strings.TrimSpace(*stateFilter)
	}

	if normalizedCity == "" && normalizedState == "" {
		return source
	}

	result := make([]photosessionmodel.PhotographerServiceAreaInterface, 0, len(source))
	for _, area := range source {
		if normalizedCity != "" && !strings.EqualFold(area.City(), normalizedCity) {
			continue
		}
		if normalizedState != "" && !strings.EqualFold(area.State(), normalizedState) {
			continue
		}
		result = append(result, area)
	}

	return result
}

func sortServiceAreas(areas []photosessionmodel.PhotographerServiceAreaInterface) {
	if len(areas) <= 1 {
		return
	}

	sort.Slice(areas, func(i, j int) bool {
		ci := strings.ToLower(strings.TrimSpace(areas[i].City()))
		cj := strings.ToLower(strings.TrimSpace(areas[j].City()))
		if ci == cj {
			si := strings.ToLower(strings.TrimSpace(areas[i].State()))
			sj := strings.ToLower(strings.TrimSpace(areas[j].State()))
			if si == sj {
				return areas[i].ID() < areas[j].ID()
			}
			return si < sj
		}
		return ci < cj
	})
}
