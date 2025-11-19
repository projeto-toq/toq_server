package complexhandler

import "github.com/gin-gonic/gin"

// ComplexHandlerPort define a interface para handlers relacionados a empreendimentos.
type ComplexHandlerPort interface {
	// GetComplexByAddress obtém detalhes do complexo a partir de um CEP e número.
	GetComplexByAddress(c *gin.Context)
}
