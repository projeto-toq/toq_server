package converters

import (
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/entity"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
)

// ToSlotModel converts SlotEntity to domain model.
func ToSlotModel(e entity.SlotEntity) photosessionmodel.PhotographerSlotInterface {
	slot := photosessionmodel.NewPhotographerSlot()
	slot.SetID(e.ID)
	slot.SetPhotographerUserID(e.PhotographerUserID)
	slot.SetSlotDate(e.SlotDate)
	slot.SetSlotStart(e.SlotStart)
	slot.SetSlotEnd(e.SlotEnd)
	slot.SetPeriod(photosessionmodel.SlotPeriod(e.Period))
	slot.SetStatus(photosessionmodel.SlotStatus(e.Status))
	slot.SetReservationToken(e.ReservationToken)
	slot.SetReservedUntil(e.ReservedUntil)
	slot.SetBookedAt(e.BookedAt)
	return slot
}
