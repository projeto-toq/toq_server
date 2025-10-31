package schedulehandler

import "github.com/gin-gonic/gin"

// ScheduleHandlerPort describes HTTP handlers responsible for agenda operations.
type ScheduleHandlerPort interface {
	GetOwnerSummary(c *gin.Context)
	GetListingAgenda(c *gin.Context)
	GetListingBlockRules(c *gin.Context)
	PostCreateBlockRule(c *gin.Context)
	PutUpdateBlockRule(c *gin.Context)
	DeleteBlockRule(c *gin.Context)
	GetListingAvailability(c *gin.Context)
	PostFinishListingAgenda(c *gin.Context)
}
