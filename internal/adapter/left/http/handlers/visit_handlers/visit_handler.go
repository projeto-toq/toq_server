package visithandlers

import (
	visithandlerport "github.com/projeto-toq/toq_server/internal/core/port/left/http/visithandler"
	visitservice "github.com/projeto-toq/toq_server/internal/core/service/visit_service"
)

// VisitHandler orchestrates HTTP endpoints for visit flows.
type VisitHandler struct {
	visitService visitservice.Service
}

// NewVisitHandler builds a new VisitHandler instance.
func NewVisitHandler(visitService visitservice.Service) visithandlerport.VisitHandlerPort {
	return &VisitHandler{visitService: visitService}
}
