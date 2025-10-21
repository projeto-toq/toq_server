package schedulehandlers

import (
	schedulehandlerport "github.com/projeto-toq/toq_server/internal/core/port/left/http/schedulehandler"
	scheduleservices "github.com/projeto-toq/toq_server/internal/core/service/schedule_service"
)

// ScheduleHandler orchestrates HTTP endpoints for property agendas.
type ScheduleHandler struct {
	scheduleService scheduleservices.ScheduleServiceInterface
}

// NewScheduleHandlerAdapter builds a new ScheduleHandler instance.
func NewScheduleHandlerAdapter(
	scheduleService scheduleservices.ScheduleServiceInterface,
) schedulehandlerport.ScheduleHandlerPort {
	return &ScheduleHandler{
		scheduleService: scheduleService,
	}
}
