package converters

import (
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/entity"
)

// ToTimeOffModel converts a time off entity into domain representation.
func ToTimeOffModel(e entity.TimeOffEntity) photosessionmodel.PhotographerTimeOffInterface {
	timeOff := photosessionmodel.NewPhotographerTimeOff()
	timeOff.SetID(e.ID)
	timeOff.SetPhotographerUserID(e.PhotographerUserID)
	timeOff.SetStartDate(e.StartDate)
	timeOff.SetEndDate(e.EndDate)
	timeOff.SetReason(e.Reason)
	return timeOff
}
