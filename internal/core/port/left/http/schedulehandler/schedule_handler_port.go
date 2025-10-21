package schedulehandler

import "github.com/gin-gonic/gin"

// ScheduleHandlerPort describes HTTP handlers responsible for agenda operations.
type ScheduleHandlerPort interface {
	PostOwnerSummary(c *gin.Context)
	PostListingAgenda(c *gin.Context)
	PostCreateBlockEntry(c *gin.Context)
	PutUpdateBlockEntry(c *gin.Context)
	DeleteBlockEntry(c *gin.Context)
	PostListAvailability(c *gin.Context)
}
