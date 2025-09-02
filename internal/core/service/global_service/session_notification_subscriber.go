package globalservice

import (
	"log/slog"
	"time"

	"github.com/giulio-alfieri/toq_server/internal/core/events"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricSessionEventHandleDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "session_event_handle_duration_seconds",
		Help:    "Time spent handling session events by type",
		Buckets: prometheus.DefBuckets,
	}, []string{"type"})
	metricDevicePruneByEvent = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "device_tokens_pruned_total",
		Help: "Device tokens pruned by event, partitioned by reason",
	}, []string{"event", "result"})
)

func init() { prometheus.MustRegister(metricSessionEventHandleDuration, metricDevicePruneByEvent) }

// StartSessionEventSubscriber subscribes to session events and triggers device-token pruning or notifications
func (gs *globalService) StartSessionEventSubscriber() func() {
	bus := gs.eventBus
	if bus == nil {
		return func() {}
	}
	unsub := bus.Subscribe(func(evt events.SessionEvent) {
		slog.Debug("session.event.received", "type", evt.Type, "user_id", evt.UserID, "session_id", evt.SessionID, "device_id", evt.DeviceID)
		start := time.Now()
		defer func() {
			metricSessionEventHandleDuration.WithLabelValues(string(evt.Type)).Observe(time.Since(start).Seconds())
		}()
		switch evt.Type {
		case events.SessionsRevoked:
			// If we have a deviceID, prune tokens associated to that device (schema fallback: no-op)
			if evt.DeviceID != "" {
				if err := gs.deviceTokenRepo.RemoveTokensByDeviceID(evt.UserID, evt.DeviceID); err != nil {
					slog.Warn("session.subscriber.device_tokens_prune_failed", "user_id", evt.UserID, "device_id", evt.DeviceID, "err", err)
					metricDevicePruneByEvent.WithLabelValues(string(evt.Type), "error").Inc()
				} else {
					slog.Info("session.subscriber.device_tokens_pruned", "user_id", evt.UserID, "device_id", evt.DeviceID)
					metricDevicePruneByEvent.WithLabelValues(string(evt.Type), "success").Inc()
				}
			}
			// Optional: send a push notification per device (requires fetching tokens by device)
			// Skipped here to avoid over-notifying; can be implemented using ListTokensByDeviceID + FCM
		}
	})
	return unsub
}
