package holidayhandlers

import (
	holidayhandlerport "github.com/projeto-toq/toq_server/internal/core/port/left/http/holidayhandler"
	holidayservices "github.com/projeto-toq/toq_server/internal/core/service/holiday_service"
)

// HolidayHandler orchestrates administrative holiday operations.
type HolidayHandler struct {
	holidayService holidayservices.HolidayServiceInterface
}

// NewHolidayHandlerAdapter builds a new holiday handler instance.
func NewHolidayHandlerAdapter(
	holidayService holidayservices.HolidayServiceInterface,
) holidayhandlerport.HolidayHandlerPort {
	return &HolidayHandler{
		holidayService: holidayService,
	}
}
