package config

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type healthStatus struct {
	Status string `json:"status"`
}

// StartHTTPHealth starts optional HTTP endpoints for health checks.
// /healthz: liveness (always 200 while process alive)
// /readyz: readiness (200 if gRPC health is SERVING, else 503)
func (c *config) StartHTTPHealth() {
	port := c.env.HEALTH.HTTPPort
	if port <= 0 {
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(healthStatus{Status: "alive"})
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if c.readiness {
			_ = json.NewEncoder(w).Encode(healthStatus{Status: "serving"})
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(healthStatus{Status: "not_serving"})
	})
	addr := ":" + strconv.Itoa(port)
	go func() {
		useTLS := c.env.HEALTH.UseTLS
		certPath := c.env.HEALTH.CertPath
		keyPath := c.env.HEALTH.KeyPath
		// Fallback to gRPC certs if not provided explicitly
		if certPath == "" || keyPath == "" {
			certPath = c.env.GRPC.CertPath
			keyPath = c.env.GRPC.KeyPath
		}
		if useTLS {
			slog.Info("HTTPS health endpoints started", "addr", addr, "cert", certPath)
			if err := http.ListenAndServeTLS(addr, certPath, keyPath, mux); err != nil {
				slog.Warn("HTTPS health endpoints failed", "error", err)
			}
			return
		}
		slog.Info("HTTP health endpoints started", "addr", addr)
		if err := http.ListenAndServe(addr, mux); err != nil {
			slog.Warn("HTTP health endpoints failed", "error", err)
		}
	}()
}
