package holidayservices

import "github.com/prometheus/client_golang/prometheus"

var (
	metricHolidayCalendarDatesDeleted = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "holiday_calendar_dates_cleaner_deleted_total",
		Help: "Total number of calendar dates deleted by the holiday cleaner",
	})
)

func init() {
	prometheus.MustRegister(metricHolidayCalendarDatesDeleted)
}
