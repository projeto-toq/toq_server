package mediaprocessingservice

import "github.com/prometheus/client_golang/prometheus"

var (
	metricMediaJobsCleanerDeleted = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "media_processing_jobs_cleaner_deleted_total",
		Help: "Total number of media processing jobs deleted by the retention cleaner",
	})
)

func init() {
	prometheus.MustRegister(metricMediaJobsCleanerDeleted)
}
