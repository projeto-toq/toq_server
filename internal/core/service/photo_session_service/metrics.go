package photosessionservices

import "github.com/prometheus/client_golang/prometheus"

var (
	metricPhotoSessionBookingsDeleted = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "photo_session_bookings_cleaner_deleted_total",
		Help: "Total number of photo session bookings deleted by the retention cleaner",
	})
	metricPhotoSessionAgendaDeleted = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "photo_session_agenda_cleaner_deleted_total",
		Help: "Total number of agenda entries deleted by the retention cleaner",
	})
)

func init() {
	prometheus.MustRegister(metricPhotoSessionBookingsDeleted)
	prometheus.MustRegister(metricPhotoSessionAgendaDeleted)
}
