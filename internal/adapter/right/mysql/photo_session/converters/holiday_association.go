package converters

import (
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/entity"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
)

// ToHolidayAssociationModel converts a DB entity to domain representation.
func ToHolidayAssociationModel(row entity.HolidayAssociation) photosessionmodel.HolidayCalendarAssociationInterface {
	model := photosessionmodel.NewHolidayCalendarAssociation()
	model.SetID(row.ID)
	model.SetPhotographerUserID(row.PhotographerUserID)
	model.SetHolidayCalendarID(row.HolidayCalendarID)
	if row.CreatedAt.Valid {
		model.SetCreatedAt(row.CreatedAt.Time)
	}
	return model
}
