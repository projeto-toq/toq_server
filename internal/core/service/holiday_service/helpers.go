package holidayservices

import (
	"strings"

	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func validateScopeInput(scope holidaymodel.CalendarScope, state, city string) error {
	switch scope {
	case holidaymodel.ScopeNational:
		return nil
	case holidaymodel.ScopeState:
		if strings.TrimSpace(state) == "" {
			return utils.ValidationError("state", "state is required for state scope")
		}
	case holidaymodel.ScopeCity:
		if strings.TrimSpace(state) == "" {
			return utils.ValidationError("state", "state is required for city scope")
		}
		if strings.TrimSpace(city) == "" {
			return utils.ValidationError("city", "city is required for city scope")
		}
	default:
		return utils.ValidationError("scope", "invalid calendar scope")
	}
	return nil
}

func cleanState(scope holidaymodel.CalendarScope, value string) string {
	if scope == holidaymodel.ScopeNational {
		return ""
	}
	state := strings.TrimSpace(value)
	if state == "" {
		return ""
	}
	return strings.ToUpper(state)
}
