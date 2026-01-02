package visitservice

import (
	"fmt"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
)

// Config stores runtime rules for visit window validation.
// MinHoursAhead defines the minimum lead time before the requested start.
// MaxDaysAhead bounds how far into the future a visit can be requested.
type Config struct {
	MinHoursAhead int
	MaxDaysAhead  int
}

// DefaultConfig returns the built-in safe defaults.
func DefaultConfig() Config {
	return Config{
		MinHoursAhead: 2,
		MaxDaysAhead:  14,
	}
}

// ConfigFromEnvironment converts YAML/env values into a Config, falling back to defaults on missing fields.
func ConfigFromEnvironment(env *globalmodel.Environment) (Config, error) {
	if env == nil {
		return DefaultConfig(), nil
	}

	cfg := Config{
		MinHoursAhead: env.Visits.MinHoursAhead,
		MaxDaysAhead:  env.Visits.MaxDaysAhead,
	}

	if cfg.MinHoursAhead <= 0 {
		cfg.MinHoursAhead = DefaultConfig().MinHoursAhead
	}
	if cfg.MaxDaysAhead <= 0 {
		cfg.MaxDaysAhead = DefaultConfig().MaxDaysAhead
	}

	if cfg.MinHoursAhead >= cfg.MaxDaysAhead*24 {
		return Config{}, fmt.Errorf("visits: min_hours_ahead must be less than max_days_ahead in hours")
	}

	return cfg, nil
}
