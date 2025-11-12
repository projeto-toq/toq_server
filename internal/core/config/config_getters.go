package config

import (
	"fmt"
	"time"

	goroutines "github.com/projeto-toq/toq_server/internal/core/go_routines"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
)

func (c *config) GetActivityTracker() *goroutines.ActivityTracker {
	return c.activityTracker
}

func (c *config) GetEnvironment() (*globalmodel.Environment, error) {
	return &c.env, nil
}

func (c *config) GetHMACSecurityConfig() (globalmodel.HMACSecurityConfig, error) {
	if c.env.SECURITY.HMAC.Secret == "" {
		return globalmodel.HMACSecurityConfig{}, fmt.Errorf("security.hmac.secret not configured")
	}
	return c.env.GetHMACSecurityConfig(), nil
}

// GetMaxWrongSigninAttempts returns the maximum number of failed signin attempts before temporary block
// Default: 3 attempts
func (c *config) GetMaxWrongSigninAttempts() int {
	if c.env.AUTH.MaxWrongSigninAttempts > 0 {
		return c.env.AUTH.MaxWrongSigninAttempts
	}
	return 3 // default fallback
}

// GetTempBlockDuration returns the duration of temporary block after max attempts exceeded
// Default: 15 minutes
func (c *config) GetTempBlockDuration() time.Duration {
	if c.env.AUTH.TempBlockDurationMinutes > 0 {
		return time.Duration(c.env.AUTH.TempBlockDurationMinutes) * time.Minute
	}
	return 15 * time.Minute // default fallback
}
