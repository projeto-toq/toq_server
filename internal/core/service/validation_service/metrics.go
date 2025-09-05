package validationservice

import "github.com/prometheus/client_golang/prometheus"

var (
	metricValidationCleanerDeleted = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "validation_cleaner_deleted_total",
		Help: "Total number of temp_user_validations rows deleted by the validation cleaner service",
	})
)

func init() {
	prometheus.MustRegister(metricValidationCleanerDeleted)
}
