package complexhandlers

import (
	complexhandlerport "github.com/projeto-toq/toq_server/internal/core/port/left/http/complexhandler"
	complexservice "github.com/projeto-toq/toq_server/internal/core/service/complex_service"
)

// ComplexHandler implementa os handlers HTTP para operações de empreendimentos publicamente acessíveis.
type ComplexHandler struct {
	complexService complexservice.ComplexServiceInterface
}

// NewComplexHandlerAdapter cria uma nova instância de ComplexHandler.
func NewComplexHandlerAdapter(
	complexService complexservice.ComplexServiceInterface,
) complexhandlerport.ComplexHandlerPort {
	return &ComplexHandler{
		complexService: complexService,
	}
}
