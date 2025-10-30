package photosessionservices

import (
	"strings"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const (
	maxServiceAreaCityLength  = 120
	maxServiceAreaStateLength = 50
)

func normalizeServiceAreaFields(city, state string) (string, string, error) {
	trimmedCity := strings.TrimSpace(city)
	if trimmedCity == "" {
		return "", "", utils.ValidationError("city", "city is required")
	}
	if len([]rune(trimmedCity)) > maxServiceAreaCityLength {
		return "", "", utils.ValidationError("city", "city must be at most 120 characters")
	}

	trimmedState := strings.TrimSpace(state)
	if trimmedState == "" {
		return "", "", utils.ValidationError("state", "state is required")
	}
	if len([]rune(trimmedState)) > maxServiceAreaStateLength {
		return "", "", utils.ValidationError("state", "state must be at most 50 characters")
	}

	return trimmedCity, strings.ToUpper(trimmedState), nil
}

func hasServiceAreaDuplicate(areas []photosessionmodel.PhotographerServiceAreaInterface, city, state string, ignoreID uint64) bool {
	for _, area := range areas {
		if ignoreID != 0 && area.ID() == ignoreID {
			continue
		}
		if strings.EqualFold(area.City(), city) && strings.EqualFold(area.State(), state) {
			return true
		}
	}
	return false
}
