package converters

import (
	"strings"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/entity"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
)

// ServiceAreaRowToModel converts a database row into a domain model instance, trimming city/state to avoid
// persisting accidental whitespace differences.
func ServiceAreaRowToModel(row entity.ServiceArea) photosessionmodel.PhotographerServiceAreaInterface {
	model := photosessionmodel.NewPhotographerServiceArea()
	model.SetID(row.ID)
	model.SetPhotographerUserID(row.PhotographerUserID)
	model.SetCity(strings.TrimSpace(row.City))
	model.SetState(strings.TrimSpace(row.State))
	return model
}

// ServiceAreaModelToRow converts a domain model into a persistence structure ready for inserts/updates.
func ServiceAreaModelToRow(model photosessionmodel.PhotographerServiceAreaInterface) entity.ServiceArea {
	return entity.ServiceArea{
		ID:                 model.ID(),
		PhotographerUserID: model.PhotographerUserID(),
		City:               model.City(),
		State:              model.State(),
	}
}
