package complexhandlers

import (
	complexhandlerport "github.com/projeto-toq/toq_server/internal/core/port/left/http/complexhandler"
	propertycoverageservice "github.com/projeto-toq/toq_server/internal/core/service/property_coverage_service"
)

// ComplexHandler implementa os handlers HTTP para operações de empreendimentos publicamente acessíveis.
type ComplexHandler struct {
	propertyCoverageService propertycoverageservice.PropertyCoverageServiceInterface
}

// NewComplexHandlerAdapter cria uma nova instância de ComplexHandler.
func NewComplexHandlerAdapter(
	propertyCoverageService propertycoverageservice.PropertyCoverageServiceInterface,
) complexhandlerport.ComplexHandlerPort {
	return &ComplexHandler{
		propertyCoverageService: propertyCoverageService,
	}
}
