package events

import "github.com/prometheus/client_golang/prometheus"

var (
	metricEventsPublishedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "event_bus_published_total",
			Help: "Total number of events published on the in-memory bus",
		},
		[]string{"type"},
	)
)

func init() {
	prometheus.MustRegister(metricEventsPublishedTotal)
}

func observePublished(t EventType) {
	metricEventsPublishedTotal.WithLabelValues(string(t)).Inc()
}
