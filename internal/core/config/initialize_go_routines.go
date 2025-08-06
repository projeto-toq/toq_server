package config

import (
	"sync"

	goroutines "github.com/giulio-alfieri/toq_server/internal/core/go_routines"
)

func (c *config) InitializeGoRoutines() {
	c.wg.Add(1)
	go goroutines.GoUpdateLastActivity(c.wg, c.userService, c.activity, c.context)
	// c.wg.Add(1)
	// go goroutines.CreciValidationWorker(c.userService, c.wg, c.context)
	c.wg.Add(1)
	go goroutines.CleanMemoryCache(&c.cache, c.wg, c.context)
}
func (c *config) GetWG() *sync.WaitGroup {
	return c.wg
}
