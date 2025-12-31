package visithandler

import "github.com/gin-gonic/gin"

// VisitHandlerPort describes HTTP handlers responsible for visit operations.
type VisitHandlerPort interface {
	CreateVisit(c *gin.Context)
	UpdateVisitStatus(c *gin.Context)
	ListVisitsRealtor(c *gin.Context)
	ListVisitsOwner(c *gin.Context)
	GetVisit(c *gin.Context)
}
