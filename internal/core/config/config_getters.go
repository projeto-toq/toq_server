package config

import (
	goroutines "github.com/giulio-alfieri/toq_server/internal/core/go_routines"
)

func (c *config) GetActivityTracker() *goroutines.ActivityTracker {
	return c.activityTracker
}
