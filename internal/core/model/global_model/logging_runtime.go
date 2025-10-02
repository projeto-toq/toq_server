package globalmodel

import (
	"log/slog"
	"strings"
	"sync"
)

// LoggingRuntimeConfig represents the runtime logging configuration provided at startup.
type LoggingRuntimeConfig struct {
	Level     slog.Level
	Format    string
	Output    string
	AddSource bool
}

var (
	loggingRuntimeConfig     LoggingRuntimeConfig
	loggingRuntimeConfigOnce sync.Once
	loggingRuntimeConfigMu   sync.RWMutex
	loggingRuntimeConfigured bool
)

// SetLoggingRuntimeConfig stores the logging runtime configuration for later consumption by adapters.
func SetLoggingRuntimeConfig(config LoggingRuntimeConfig) {
	config.Format = strings.ToLower(config.Format)
	config.Output = strings.ToLower(config.Output)

	loggingRuntimeConfigOnce.Do(func() {
		loggingRuntimeConfigured = true
	})

	loggingRuntimeConfigMu.Lock()
	loggingRuntimeConfig = config
	loggingRuntimeConfigMu.Unlock()
}

// GetLoggingRuntimeConfig returns the logging runtime configuration if it was previously defined.
func GetLoggingRuntimeConfig() (LoggingRuntimeConfig, bool) {
	loggingRuntimeConfigMu.RLock()
	defer loggingRuntimeConfigMu.RUnlock()

	if !loggingRuntimeConfigured {
		return LoggingRuntimeConfig{}, false
	}

	return loggingRuntimeConfig, true
}
