package config

import (
	"sync"
	"time"

	goroutines "github.com/giulio-alfieri/toq_server/internal/core/go_routines"
)

func (c *config) InitializeGoRoutines() {
	// Start Activity Tracker worker
	c.wg.Add(1)
	go c.activityTracker.StartBatchWorker(c.wg, c.context)

	c.wg.Add(1)
	go goroutines.CreciValidationWorker(c.userService, c.wg, c.context)

	c.wg.Add(1)
	go goroutines.CleanMemoryCache(&c.cache, c.wg, c.context)

	// Start Session Cleaner worker (default interval 5m if not configured)
	intervalSeconds := c.env.AUTH.SessionCleanerIntervalSeconds
	if intervalSeconds <= 0 {
		intervalSeconds = 300 // 5 minutes default
	}
	c.wg.Add(1)
	go goroutines.SessionCleaner(c.db, c.wg, c.context, time.Duration(intervalSeconds)*time.Second)
}
func (c *config) GetWG() *sync.WaitGroup {
	return c.wg
}
