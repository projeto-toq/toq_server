package sessionservice

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricSessionCleanerDeleted = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "session_cleaner_deleted_total",
		Help: "Total number of sessions deleted by the session cleaner service",
	})
)

func init() {
	prometheus.MustRegister(metricSessionCleanerDeleted)
}
