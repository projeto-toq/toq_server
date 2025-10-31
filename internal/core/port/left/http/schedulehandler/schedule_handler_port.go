package schedulehandler

import "github.com/gin-gonic/gin"

// ScheduleHandlerPort describes HTTP handlers responsible for agenda operations.
type ScheduleHandlerPort interface {
	GetOwnerSummary(c *gin.Context)
	GetListingAgenda(c *gin.Context)
	GetListingBlockEntries(c *gin.Context)
	PostCreateBlockEntry(c *gin.Context)
	PutUpdateBlockEntry(c *gin.Context)
	DeleteBlockEntry(c *gin.Context)
	GetListingAvailability(c *gin.Context)
	PostFinishListingAgenda(c *gin.Context)
}
