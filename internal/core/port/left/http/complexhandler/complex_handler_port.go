package complexhandler

import "github.com/gin-gonic/gin"

// ComplexHandlerPort define a interface para handlers relacionados a empreendimentos.
type ComplexHandlerPort interface {
	// ListSizesByAddress lista tamanhos disponíveis a partir de um CEP e número.
	ListSizesByAddress(c *gin.Context)
}
