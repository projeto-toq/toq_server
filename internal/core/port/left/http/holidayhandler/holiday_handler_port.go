package holidayhandler

import "github.com/gin-gonic/gin"

// HolidayHandlerPort exposes administrative handlers for holiday calendars and dates.
type HolidayHandlerPort interface {
	ListCalendars(c *gin.Context)
	GetCalendarDetail(c *gin.Context)
	CreateCalendar(c *gin.Context)
	UpdateCalendar(c *gin.Context)
	CreateCalendarDate(c *gin.Context)
	ListCalendarDates(c *gin.Context)
	DeleteCalendarDate(c *gin.Context)
}
