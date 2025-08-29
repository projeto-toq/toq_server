package config

import (
	goroutines "github.com/giulio-alfieri/toq_server/internal/core/go_routines"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
)

func (c *config) GetActivityTracker() *goroutines.ActivityTracker {
	return c.activityTracker
}

func (c *config) GetEnvironment() (*globalmodel.Environment, error) {
	return &c.env, nil
}
